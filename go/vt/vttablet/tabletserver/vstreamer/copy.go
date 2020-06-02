/*
Copyright 2020 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
TestVStreamCopyCompleteFlow tests a complete happy VStream Copy flow: copy/catchup/fastforward/replicate
Three tables t1, t2, t3 are copied. Initially 10 (numInitialRows) rows are inserted into each.
To avoid races in testing we send additional events when *uvstreamerTestMode* is set to true. These are used in
conjunction with callbacks to do additional crud at precise points of the flow to test the different paths
We intercept the vstreamer send callback to look for specific events and invoke these test callbacks.
Fast forward requires tables to be locked briefly to get a snapshot: the test uses this knowledge to hold a lock
on the table in order to insert rows for fastforward to find.

The flow is as follows:
	t1: copy phase, 10 rows.
	t2: copy phase to start. Test event is sent, intercepted and a row in inserted into t1
    t1: fastforward finds this event

*/

package vstreamer

import (
	"context"
	"fmt"
	"io"
	"math"
	"time"

	"vitess.io/vitess/go/mysql"
	"vitess.io/vitess/go/sqltypes"
	"vitess.io/vitess/go/vt/log"
	binlogdatapb "vitess.io/vitess/go/vt/proto/binlogdata"
	querypb "vitess.io/vitess/go/vt/proto/query"
)

func (uvs *uvstreamer) copy(ctx context.Context) error {
	for len(uvs.tablesToCopy) > 0 {
		tableName := uvs.tablesToCopy[0]
		log.Infof("Copystate not empty starting catchupAndCopy on table %s", tableName)
		if err := uvs.catchupAndCopy(ctx, tableName); err != nil {
			return err
		}
		//time.Sleep(2*time.Second) //FIXME for debugging
	}
	log.Info("No tables left to copy")
	return nil
}

func (uvs *uvstreamer) catchupAndCopy(ctx context.Context, tableName string) error {
	log.Infof("catchupAndCopy for %s", tableName)
	if !uvs.pos.IsZero() {
		if err := uvs.catchup(ctx); err != nil {
			log.Infof("catchupAndCopy: catchup returned %v", err)
			return err
		}
	}

	log.Infof("catchupAndCopy: before copyTable %s", tableName)
	uvs.fields = nil
	return uvs.copyTable(ctx, tableName)
}

func (uvs *uvstreamer) catchup(ctx context.Context) error {
	log.Infof("starting catchup ...")
	uvs.setSecondsBehindMaster(math.MaxInt64)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errch := make(chan error, 1)
	go func() {
		startPos := mysql.EncodePosition(uvs.pos)
		vs := newVStreamer(ctx, uvs.cp, uvs.se, uvs.sh, startPos, "", uvs.filter, uvs.vschema, uvs.send2)
		errch <- vs.Stream()
		uvs.vs = nil
		log.Infof("catchup vs.stream returned with vs.pos %s", vs.pos.String())
	}()

	// Wait for catchup.
	tkr := time.NewTicker(uvs.config.CatchupRetryTime)
	defer tkr.Stop()
	seconds := int64(uvs.config.MaxReplicationLag / time.Second)
	for {
		sbm := uvs.getSecondsBehindMaster()
		log.Infof("Checking sbm %d vs config %d", sbm, seconds)
		if sbm <= seconds {
			log.Infof("Canceling context because lag is %d:%d", sbm, seconds)
			cancel()
			// Make sure vplayer returns before returning.
			<-errch
			return nil
		}
		select {
		case err := <-errch:
			if err != nil {
				return err
			}
			return io.EOF
		case <-ctx.Done():
			// Make sure vplayer returns before returning.
			<-errch
			return io.EOF
		case <-tkr.C:
		}
	}
}

func (uvs *uvstreamer) sendFieldEvent(ctx context.Context, gtid string, fieldEvent *binlogdatapb.FieldEvent) error {
	evs := []*binlogdatapb.VEvent{{
		Type: binlogdatapb.VEventType_BEGIN,
	}, {
		Type:       binlogdatapb.VEventType_FIELD,
		FieldEvent: fieldEvent,
	}}
	log.Infof("Sending field event %v, gtid is %s", fieldEvent, gtid)
	uvs.send(evs)
	if err := uvs.setPosition(gtid, true); err != nil {
		log.Infof("setPosition returned error %v", err)
		return err
	}
	return nil

}

// send one RowEvent per row, followed by a LastPK (merged in VTGate with vgtid)
func (uvs *uvstreamer) sendEventsForRows(ctx context.Context, tableName string, rows *binlogdatapb.VStreamRowsResponse, qr *querypb.QueryResult) error {
	var evs []*binlogdatapb.VEvent
	for _, row := range rows.Rows {
		ev := &binlogdatapb.VEvent{
			Type: binlogdatapb.VEventType_ROW,
			RowEvent: &binlogdatapb.RowEvent{
				TableName: tableName,
				RowChanges: []*binlogdatapb.RowChange{{
					Before: nil,
					After:  row,
				}},
			},
		}
		evs = append(evs, ev)
	}
	lastPKEvent := &binlogdatapb.LastPKEvent{
		TableLastPK: &binlogdatapb.TableLastPK{
			TableName: tableName,
			Lastpk:    qr,
		},
		Completed: false,
	}

	ev := &binlogdatapb.VEvent{
		Type:        binlogdatapb.VEventType_LASTPK,
		LastPKEvent: lastPKEvent,
	}
	evs = append(evs, ev)
	evs = append(evs, &binlogdatapb.VEvent{
		Type: binlogdatapb.VEventType_COMMIT,
	})
	if err := uvs.send(evs); err != nil {
		log.Infof("send returned error %v", err)
		return err
	}
	return nil
}

func getLastPKFromQR(qr *querypb.QueryResult) []sqltypes.Value {
	var lastPK []sqltypes.Value
	r := sqltypes.Proto3ToResult(qr)
	if len(r.Rows) != 1 {
		log.Errorf("unexpected lastpk input: %v", qr)
		return nil
	}
	lastPK = r.Rows[0]
	return lastPK
}

func getQRFromLastPK(fields []*querypb.Field, lastPK []sqltypes.Value) *querypb.QueryResult {
	row := sqltypes.RowToProto3(lastPK)
	qr := &querypb.QueryResult{
		Fields: fields,
		Rows:   []*querypb.Row{row},
	}
	return qr
}

func (uvs *uvstreamer) copyTable(ctx context.Context, tableName string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var newLastPK *sqltypes.Result
	lastPK := getLastPKFromQR(uvs.plans[tableName].tablePK.Lastpk)
	filter := uvs.plans[tableName].rule.Filter

	log.Infof("Starting copyTable for %s, PK %v", tableName, lastPK)
	uvs.sendTestEvent(fmt.Sprintf("Copy Start %s", tableName))

	err := uvs.vse.StreamRows(ctx, filter, lastPK, func(rows *binlogdatapb.VStreamRowsResponse) error {
		select {
		case <-ctx.Done():
			log.Infof("Returning io.EOF in StreamRows")
			return io.EOF
		default:
		}
		if uvs.fields == nil {
			if len(rows.Fields) == 0 {
				return fmt.Errorf("expecting field event first, got: %v", rows)
			}
			pos, _ := mysql.DecodePosition(rows.Gtid)
			if !uvs.pos.IsZero() && !uvs.pos.AtLeast(pos) {
				if err := uvs.fastForward(rows.Gtid); err != nil {
					log.Infof("fastForward returned error %v", err)
					return err
				}
				if mysql.EncodePosition(uvs.pos) != rows.Gtid {
					return fmt.Errorf("position after fastforward was %s but stopPos was %s", uvs.pos, rows.Gtid)
				}
				if err := uvs.setPosition(rows.Gtid, false); err != nil {
					return err
				}
			} else {
				log.Infof("Not starting fastforward pos is %s, uvs.pos is %s, rows.gtid %s", pos, uvs.pos, rows.Gtid)
			}

			fieldEvent := &binlogdatapb.FieldEvent{
				TableName: tableName,
				Fields:    rows.Fields,
			}
			uvs.fields = rows.Fields
			uvs.pkfields = rows.Pkfields
			if err := uvs.sendFieldEvent(ctx, rows.Gtid, fieldEvent); err != nil {
				log.Infof("sendFieldEvent returned error %v", err)
				return err
			}
		}
		if len(rows.Rows) == 0 {
			log.Infof("0 rows returned for table %s", tableName)
			return nil
		}

		newLastPK = sqltypes.CustomProto3ToResult(uvs.pkfields, &querypb.QueryResult{
			Fields: rows.Fields,
			Rows:   []*querypb.Row{rows.Lastpk},
		})
		qrLastPK := sqltypes.ResultToProto3(newLastPK)
		log.Infof("Calling sendEventForRows with gtid %s", rows.Gtid)
		if err := uvs.sendEventsForRows(ctx, tableName, rows, qrLastPK); err != nil {
			log.Infof("sendEventsForRows returned error %v", err)
			return err
		}

		uvs.setCopyState(tableName, qrLastPK)
		log.Infof("NewLastPK: %v", qrLastPK)
		return nil
	})
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		log.Infof("Context done: Copy of %v stopped at lastpk: %v", tableName, newLastPK)
		return ctx.Err()
	default:
	}

	log.Infof("Copy of %v finished at lastpk: %v", tableName, newLastPK)
	if err := uvs.copyComplete(tableName); err != nil {
		return err
	}
	return nil
}

func (uvs *uvstreamer) fastForward(stopPos string) error {
	log.Infof("starting fastForward from %s upto pos %s", mysql.EncodePosition(uvs.pos), stopPos)
	uvs.stopPos, _ = mysql.DecodePosition(stopPos)
	vs := newVStreamer(uvs.ctx, uvs.cp, uvs.se, uvs.sh, mysql.EncodePosition(uvs.pos), "", uvs.filter, uvs.vschema, uvs.send2)
	uvs.vs = vs
	return vs.Stream()
}
