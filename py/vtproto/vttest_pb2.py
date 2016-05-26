# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: vttest.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
from google.protobuf import descriptor_pb2
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='vttest.proto',
  package='vttest',
  syntax='proto3',
  serialized_pb=_b('\n\x0cvttest.proto\x12\x06vttest\"/\n\x05Shard\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x18\n\x10\x64\x62_name_override\x18\x02 \x01(\t\"\x88\x01\n\x08Keyspace\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x1d\n\x06shards\x18\x02 \x03(\x0b\x32\r.vttest.Shard\x12\x1c\n\x14sharding_column_name\x18\x03 \x01(\t\x12\x1c\n\x14sharding_column_type\x18\x04 \x01(\t\x12\x13\n\x0bserved_from\x18\x05 \x01(\t\"5\n\x0eVTTestTopology\x12#\n\tkeyspaces\x18\x01 \x03(\x0b\x32\x10.vttest.Keyspaceb\x06proto3')
)
_sym_db.RegisterFileDescriptor(DESCRIPTOR)




_SHARD = _descriptor.Descriptor(
  name='Shard',
  full_name='vttest.Shard',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='vttest.Shard.name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='db_name_override', full_name='vttest.Shard.db_name_override', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=24,
  serialized_end=71,
)


_KEYSPACE = _descriptor.Descriptor(
  name='Keyspace',
  full_name='vttest.Keyspace',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='vttest.Keyspace.name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='shards', full_name='vttest.Keyspace.shards', index=1,
      number=2, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='sharding_column_name', full_name='vttest.Keyspace.sharding_column_name', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='sharding_column_type', full_name='vttest.Keyspace.sharding_column_type', index=3,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='served_from', full_name='vttest.Keyspace.served_from', index=4,
      number=5, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=74,
  serialized_end=210,
)


_VTTESTTOPOLOGY = _descriptor.Descriptor(
  name='VTTestTopology',
  full_name='vttest.VTTestTopology',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='keyspaces', full_name='vttest.VTTestTopology.keyspaces', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=212,
  serialized_end=265,
)

_KEYSPACE.fields_by_name['shards'].message_type = _SHARD
_VTTESTTOPOLOGY.fields_by_name['keyspaces'].message_type = _KEYSPACE
DESCRIPTOR.message_types_by_name['Shard'] = _SHARD
DESCRIPTOR.message_types_by_name['Keyspace'] = _KEYSPACE
DESCRIPTOR.message_types_by_name['VTTestTopology'] = _VTTESTTOPOLOGY

Shard = _reflection.GeneratedProtocolMessageType('Shard', (_message.Message,), dict(
  DESCRIPTOR = _SHARD,
  __module__ = 'vttest_pb2'
  # @@protoc_insertion_point(class_scope:vttest.Shard)
  ))
_sym_db.RegisterMessage(Shard)

Keyspace = _reflection.GeneratedProtocolMessageType('Keyspace', (_message.Message,), dict(
  DESCRIPTOR = _KEYSPACE,
  __module__ = 'vttest_pb2'
  # @@protoc_insertion_point(class_scope:vttest.Keyspace)
  ))
_sym_db.RegisterMessage(Keyspace)

VTTestTopology = _reflection.GeneratedProtocolMessageType('VTTestTopology', (_message.Message,), dict(
  DESCRIPTOR = _VTTESTTOPOLOGY,
  __module__ = 'vttest_pb2'
  # @@protoc_insertion_point(class_scope:vttest.VTTestTopology)
  ))
_sym_db.RegisterMessage(VTTestTopology)


import abc
from grpc.beta import implementations as beta_implementations
from grpc.framework.common import cardinality
from grpc.framework.interfaces.face import utilities as face_utilities
# @@protoc_insertion_point(module_scope)
