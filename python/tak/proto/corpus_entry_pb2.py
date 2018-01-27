# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: tak/proto/corpus_entry.proto

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
  name='tak/proto/corpus_entry.proto',
  package='tak.proto',
  syntax='proto3',
  serialized_pb=_b('\n\x1ctak/proto/corpus_entry.proto\x12\ttak.proto\"~\n\x0b\x43orpusEntry\x12\x0b\n\x03\x64\x61y\x18\x01 \x01(\t\x12\n\n\x02id\x18\x02 \x01(\x05\x12\x0b\n\x03ply\x18\x03 \x01(\x05\x12\x0b\n\x03tps\x18\x04 \x01(\t\x12\x0c\n\x04move\x18\x05 \x01(\t\x12\r\n\x05value\x18\x06 \x01(\x02\x12\r\n\x05plies\x18\x07 \x01(\x05\x12\x10\n\x08\x66\x65\x61tures\x18\x08 \x03(\x03\x42\x04Z\x02pbb\x06proto3')
)




_CORPUSENTRY = _descriptor.Descriptor(
  name='CorpusEntry',
  full_name='tak.proto.CorpusEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='day', full_name='tak.proto.CorpusEntry.day', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='id', full_name='tak.proto.CorpusEntry.id', index=1,
      number=2, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='ply', full_name='tak.proto.CorpusEntry.ply', index=2,
      number=3, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='tps', full_name='tak.proto.CorpusEntry.tps', index=3,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='move', full_name='tak.proto.CorpusEntry.move', index=4,
      number=5, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='value', full_name='tak.proto.CorpusEntry.value', index=5,
      number=6, type=2, cpp_type=6, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='plies', full_name='tak.proto.CorpusEntry.plies', index=6,
      number=7, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='features', full_name='tak.proto.CorpusEntry.features', index=7,
      number=8, type=3, cpp_type=2, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
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
  serialized_start=43,
  serialized_end=169,
)

DESCRIPTOR.message_types_by_name['CorpusEntry'] = _CORPUSENTRY
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

CorpusEntry = _reflection.GeneratedProtocolMessageType('CorpusEntry', (_message.Message,), dict(
  DESCRIPTOR = _CORPUSENTRY,
  __module__ = 'tak.proto.corpus_entry_pb2'
  # @@protoc_insertion_point(class_scope:tak.proto.CorpusEntry)
  ))
_sym_db.RegisterMessage(CorpusEntry)


DESCRIPTOR.has_options = True
DESCRIPTOR._options = _descriptor._ParseOptions(descriptor_pb2.FileOptions(), _b('Z\002pb'))
# @@protoc_insertion_point(module_scope)
