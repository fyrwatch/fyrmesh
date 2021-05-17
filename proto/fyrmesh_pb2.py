# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: proto/fyrmesh.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='proto/fyrmesh.proto',
  package='main',
  syntax='proto3',
  serialized_options=b'Z\006/proto',
  create_key=_descriptor._internal_create_key,
  serialized_pb=b'\n\x13proto/fyrmesh.proto\x12\x04main\"\x81\x01\n\x07Trigger\x12\x16\n\x0etriggermessage\x18\x01 \x01(\t\x12-\n\x08metadata\x18\x02 \x03(\x0b\x32\x1b.main.Trigger.MetadataEntry\x1a/\n\rMetadataEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\"-\n\x0b\x41\x63knowledge\x12\x0f\n\x07success\x18\x01 \x01(\x08\x12\r\n\x05\x65rror\x18\x02 \x01(\t\"/\n\nMeshStatus\x12\x0e\n\x06meshID\x18\x01 \x01(\t\x12\x11\n\tconnected\x18\x02 \x01(\x08\"\x1c\n\tSimpleLog\x12\x0f\n\x07message\x18\x01 \x01(\t\"U\n\nComplexLog\x12\x11\n\tlogsource\x18\x01 \x01(\t\x12\x0f\n\x07logtype\x18\x02 \x01(\t\x12\x0f\n\x07logtime\x18\x03 \x01(\t\x12\x12\n\nlogmessage\x18\x04 \x01(\t\"\x88\x01\n\x0e\x43ontrolCommand\x12\x0f\n\x07\x63ommand\x18\x01 \x01(\t\x12\x34\n\x08metadata\x18\x02 \x03(\x0b\x32\".main.ControlCommand.MetadataEntry\x1a/\n\rMetadataEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\x32l\n\tInterface\x12+\n\x04Read\x12\r.main.Trigger\x1a\x10.main.ComplexLog\"\x00\x30\x01\x12\x32\n\x05Write\x12\x14.main.ControlCommand\x1a\x11.main.Acknowledge\"\x00\x32\xc8\x01\n\x0cOrchestrator\x12+\n\x06Status\x12\r.main.Trigger\x1a\x10.main.MeshStatus\"\x00\x12\x30\n\nConnection\x12\r.main.Trigger\x1a\x11.main.Acknowledge\"\x00\x12-\n\x07Observe\x12\r.main.Trigger\x1a\x0f.main.SimpleLog\"\x00\x30\x01\x12*\n\x04Ping\x12\r.main.Trigger\x1a\x11.main.Acknowledge\"\x00\x42\x08Z\x06/protob\x06proto3'
)




_TRIGGER_METADATAENTRY = _descriptor.Descriptor(
  name='MetadataEntry',
  full_name='main.Trigger.MetadataEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='main.Trigger.MetadataEntry.key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='value', full_name='main.Trigger.MetadataEntry.value', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=b'8\001',
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=112,
  serialized_end=159,
)

_TRIGGER = _descriptor.Descriptor(
  name='Trigger',
  full_name='main.Trigger',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='triggermessage', full_name='main.Trigger.triggermessage', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='metadata', full_name='main.Trigger.metadata', index=1,
      number=2, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[_TRIGGER_METADATAENTRY, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=30,
  serialized_end=159,
)


_ACKNOWLEDGE = _descriptor.Descriptor(
  name='Acknowledge',
  full_name='main.Acknowledge',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='success', full_name='main.Acknowledge.success', index=0,
      number=1, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='error', full_name='main.Acknowledge.error', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=161,
  serialized_end=206,
)


_MESHSTATUS = _descriptor.Descriptor(
  name='MeshStatus',
  full_name='main.MeshStatus',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='meshID', full_name='main.MeshStatus.meshID', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='connected', full_name='main.MeshStatus.connected', index=1,
      number=2, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=208,
  serialized_end=255,
)


_SIMPLELOG = _descriptor.Descriptor(
  name='SimpleLog',
  full_name='main.SimpleLog',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='message', full_name='main.SimpleLog.message', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=257,
  serialized_end=285,
)


_COMPLEXLOG = _descriptor.Descriptor(
  name='ComplexLog',
  full_name='main.ComplexLog',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='logsource', full_name='main.ComplexLog.logsource', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='logtype', full_name='main.ComplexLog.logtype', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='logtime', full_name='main.ComplexLog.logtime', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='logmessage', full_name='main.ComplexLog.logmessage', index=3,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=287,
  serialized_end=372,
)


_CONTROLCOMMAND_METADATAENTRY = _descriptor.Descriptor(
  name='MetadataEntry',
  full_name='main.ControlCommand.MetadataEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='main.ControlCommand.MetadataEntry.key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='value', full_name='main.ControlCommand.MetadataEntry.value', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=b'8\001',
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=112,
  serialized_end=159,
)

_CONTROLCOMMAND = _descriptor.Descriptor(
  name='ControlCommand',
  full_name='main.ControlCommand',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='command', full_name='main.ControlCommand.command', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='metadata', full_name='main.ControlCommand.metadata', index=1,
      number=2, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[_CONTROLCOMMAND_METADATAENTRY, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=375,
  serialized_end=511,
)

_TRIGGER_METADATAENTRY.containing_type = _TRIGGER
_TRIGGER.fields_by_name['metadata'].message_type = _TRIGGER_METADATAENTRY
_CONTROLCOMMAND_METADATAENTRY.containing_type = _CONTROLCOMMAND
_CONTROLCOMMAND.fields_by_name['metadata'].message_type = _CONTROLCOMMAND_METADATAENTRY
DESCRIPTOR.message_types_by_name['Trigger'] = _TRIGGER
DESCRIPTOR.message_types_by_name['Acknowledge'] = _ACKNOWLEDGE
DESCRIPTOR.message_types_by_name['MeshStatus'] = _MESHSTATUS
DESCRIPTOR.message_types_by_name['SimpleLog'] = _SIMPLELOG
DESCRIPTOR.message_types_by_name['ComplexLog'] = _COMPLEXLOG
DESCRIPTOR.message_types_by_name['ControlCommand'] = _CONTROLCOMMAND
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

Trigger = _reflection.GeneratedProtocolMessageType('Trigger', (_message.Message,), {

  'MetadataEntry' : _reflection.GeneratedProtocolMessageType('MetadataEntry', (_message.Message,), {
    'DESCRIPTOR' : _TRIGGER_METADATAENTRY,
    '__module__' : 'proto.fyrmesh_pb2'
    # @@protoc_insertion_point(class_scope:main.Trigger.MetadataEntry)
    })
  ,
  'DESCRIPTOR' : _TRIGGER,
  '__module__' : 'proto.fyrmesh_pb2'
  # @@protoc_insertion_point(class_scope:main.Trigger)
  })
_sym_db.RegisterMessage(Trigger)
_sym_db.RegisterMessage(Trigger.MetadataEntry)

Acknowledge = _reflection.GeneratedProtocolMessageType('Acknowledge', (_message.Message,), {
  'DESCRIPTOR' : _ACKNOWLEDGE,
  '__module__' : 'proto.fyrmesh_pb2'
  # @@protoc_insertion_point(class_scope:main.Acknowledge)
  })
_sym_db.RegisterMessage(Acknowledge)

MeshStatus = _reflection.GeneratedProtocolMessageType('MeshStatus', (_message.Message,), {
  'DESCRIPTOR' : _MESHSTATUS,
  '__module__' : 'proto.fyrmesh_pb2'
  # @@protoc_insertion_point(class_scope:main.MeshStatus)
  })
_sym_db.RegisterMessage(MeshStatus)

SimpleLog = _reflection.GeneratedProtocolMessageType('SimpleLog', (_message.Message,), {
  'DESCRIPTOR' : _SIMPLELOG,
  '__module__' : 'proto.fyrmesh_pb2'
  # @@protoc_insertion_point(class_scope:main.SimpleLog)
  })
_sym_db.RegisterMessage(SimpleLog)

ComplexLog = _reflection.GeneratedProtocolMessageType('ComplexLog', (_message.Message,), {
  'DESCRIPTOR' : _COMPLEXLOG,
  '__module__' : 'proto.fyrmesh_pb2'
  # @@protoc_insertion_point(class_scope:main.ComplexLog)
  })
_sym_db.RegisterMessage(ComplexLog)

ControlCommand = _reflection.GeneratedProtocolMessageType('ControlCommand', (_message.Message,), {

  'MetadataEntry' : _reflection.GeneratedProtocolMessageType('MetadataEntry', (_message.Message,), {
    'DESCRIPTOR' : _CONTROLCOMMAND_METADATAENTRY,
    '__module__' : 'proto.fyrmesh_pb2'
    # @@protoc_insertion_point(class_scope:main.ControlCommand.MetadataEntry)
    })
  ,
  'DESCRIPTOR' : _CONTROLCOMMAND,
  '__module__' : 'proto.fyrmesh_pb2'
  # @@protoc_insertion_point(class_scope:main.ControlCommand)
  })
_sym_db.RegisterMessage(ControlCommand)
_sym_db.RegisterMessage(ControlCommand.MetadataEntry)


DESCRIPTOR._options = None
_TRIGGER_METADATAENTRY._options = None
_CONTROLCOMMAND_METADATAENTRY._options = None

_INTERFACE = _descriptor.ServiceDescriptor(
  name='Interface',
  full_name='main.Interface',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  create_key=_descriptor._internal_create_key,
  serialized_start=513,
  serialized_end=621,
  methods=[
  _descriptor.MethodDescriptor(
    name='Read',
    full_name='main.Interface.Read',
    index=0,
    containing_service=None,
    input_type=_TRIGGER,
    output_type=_COMPLEXLOG,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='Write',
    full_name='main.Interface.Write',
    index=1,
    containing_service=None,
    input_type=_CONTROLCOMMAND,
    output_type=_ACKNOWLEDGE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
])
_sym_db.RegisterServiceDescriptor(_INTERFACE)

DESCRIPTOR.services_by_name['Interface'] = _INTERFACE


_ORCHESTRATOR = _descriptor.ServiceDescriptor(
  name='Orchestrator',
  full_name='main.Orchestrator',
  file=DESCRIPTOR,
  index=1,
  serialized_options=None,
  create_key=_descriptor._internal_create_key,
  serialized_start=624,
  serialized_end=824,
  methods=[
  _descriptor.MethodDescriptor(
    name='Status',
    full_name='main.Orchestrator.Status',
    index=0,
    containing_service=None,
    input_type=_TRIGGER,
    output_type=_MESHSTATUS,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='Connection',
    full_name='main.Orchestrator.Connection',
    index=1,
    containing_service=None,
    input_type=_TRIGGER,
    output_type=_ACKNOWLEDGE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='Observe',
    full_name='main.Orchestrator.Observe',
    index=2,
    containing_service=None,
    input_type=_TRIGGER,
    output_type=_SIMPLELOG,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='Ping',
    full_name='main.Orchestrator.Ping',
    index=3,
    containing_service=None,
    input_type=_TRIGGER,
    output_type=_ACKNOWLEDGE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
])
_sym_db.RegisterServiceDescriptor(_ORCHESTRATOR)

DESCRIPTOR.services_by_name['Orchestrator'] = _ORCHESTRATOR

# @@protoc_insertion_point(module_scope)
