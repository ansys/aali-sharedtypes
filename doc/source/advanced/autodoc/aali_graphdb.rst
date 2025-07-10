.. _aali_graphdb:

aali_graphdb
============

This package provides functionality for aali_graphdb.

Functions
---------

.. function:: TestLogicalTypesMarshal

   func TestLogicalTypesMarshal(t \*testing.T)


.. function:: TestValuesMarshal

   func TestValuesMarshal(t \*testing.T)


.. function:: NewClient

   func NewClient(address string, httpClient \*http.Client) (\*Client, error)


.. function:: DefaultClient

   func DefaultClient(address string) (\*Client, error)


.. function:: TestGetHealth

   func TestGetHealth(t \*testing.T)


.. function:: TestGetDatabases

   func TestGetDatabases(t \*testing.T)


.. function:: TestCreateDatabase

   func TestCreateDatabase(t \*testing.T)


.. function:: TestDeleteDatabase

   func TestDeleteDatabase(t \*testing.T)


.. function:: TestReadWriteData

   func TestReadWriteData(t \*testing.T)


.. function:: TestReadWriteGeneric

   func TestReadWriteGeneric(t \*testing.T)


.. function:: TestReadWriteDataWithParameters

   func TestReadWriteDataWithParameters(t \*testing.T)


.. function:: TestParametersStruct

   func TestParametersStruct(t \*testing.T)


.. function:: TestErrorsReturned

   func TestErrorsReturned(t \*testing.T)


Types
-----

.. type:: LogicalType

   type LogicalType interface


.. type:: AnyLogicalType

   type AnyLogicalType struct


.. type:: BoolLogicalType

   type BoolLogicalType struct


.. type:: SerialLogicalType

   type SerialLogicalType struct


.. type:: Int64LogicalType

   type Int64LogicalType struct


.. type:: Int32LogicalType

   type Int32LogicalType struct


.. type:: Int16LogicalType

   type Int16LogicalType struct


.. type:: Int8LogicalType

   type Int8LogicalType struct


.. type:: UInt64LogicalType

   type UInt64LogicalType struct


.. type:: UInt32LogicalType

   type UInt32LogicalType struct


.. type:: UInt16LogicalType

   type UInt16LogicalType struct


.. type:: UInt8LogicalType

   type UInt8LogicalType struct


.. type:: Int128LogicalType

   type Int128LogicalType struct


.. type:: DoubleLogicalType

   type DoubleLogicalType struct


.. type:: FloatLogicalType

   type FloatLogicalType struct


.. type:: DateLogicalType

   type DateLogicalType struct


.. type:: IntervalLogicalType

   type IntervalLogicalType struct


.. type:: TimestampLogicalType

   type TimestampLogicalType struct


.. type:: TimestampTzLogicalType

   type TimestampTzLogicalType struct


.. type:: TimestampNsLogicalType

   type TimestampNsLogicalType struct


.. type:: TimestampMsLogicalType

   type TimestampMsLogicalType struct


.. type:: TimestampSecLogicalType

   type TimestampSecLogicalType struct


.. type:: InternalIDTypeLogicalType

   type InternalIDTypeLogicalType struct


.. type:: StringLogicalType

   type StringLogicalType struct


.. type:: BlobLogicalType

   type BlobLogicalType struct


.. type:: ListLogicalType

   type ListLogicalType struct


.. type:: ArrayLogicalType

   type ArrayLogicalType struct


.. type:: StructLogicalType

   type StructLogicalType struct


.. type:: NodeLogicalType

   type NodeLogicalType struct


.. type:: RelLogicalType

   type RelLogicalType struct


.. type:: RecursiveRelLogicalType

   type RecursiveRelLogicalType struct


.. type:: MapLogicalType

   type MapLogicalType struct


.. type:: UnionLogicalType

   type UnionLogicalType struct


.. type:: UUIDLogicalType

   type UUIDLogicalType struct


.. type:: DecimalLogicalType

   type DecimalLogicalType struct


.. type:: Client

   type Client struct


.. type:: ParameterMap

   type ParameterMap map\[string\]Value


.. type:: Parameters

   type Parameters interface


.. type:: StdoutLogConsumer

   type StdoutLogConsumer struct

   StdoutLogConsumer is a LogConsumer that prints the log to stdout


.. type:: ParamsStruct

   type ParamsStruct struct


.. type:: Value

   type Value interface


.. type:: NullValue

   type NullValue struct


.. type:: BoolValue

   type BoolValue bool


.. type:: Int64Value

   type Int64Value int64


.. type:: Int32Value

   type Int32Value int32


.. type:: Int16Value

   type Int16Value int16


.. type:: Int8Value

   type Int8Value int8


.. type:: UInt64Value

   type UInt64Value uint64


.. type:: UInt32Value

   type UInt32Value uint32


.. type:: UInt16Value

   type UInt16Value uint16


.. type:: UInt8Value

   type UInt8Value uint8


.. type:: Int128Value

   type Int128Value int64 // no int128 in go, but could still be that type in the DB


.. type:: DoubleValue

   type DoubleValue float64


.. type:: FloatValue

   type FloatValue float32


.. type:: DateValue

   type DateValue civil.Date


.. type:: IntervalValue

   type IntervalValue time.Duration


.. type:: TimestampValue

   type TimestampValue time.Time


.. type:: TimestampTzValue

   type TimestampTzValue time.Time


.. type:: TimestampNsValue

   type TimestampNsValue time.Time


.. type:: TimestampMsValue

   type TimestampMsValue time.Time


.. type:: TimestampSecValue

   type TimestampSecValue time.Time


.. type:: InternalID

   type InternalID struct


.. type:: InternalIDValue

   type InternalIDValue struct


.. type:: StringValue

   type StringValue string


.. type:: BlobValue

   type BlobValue \[\]uint8


.. type:: ListValue

   type ListValue struct


.. type:: ArrayValue

   type ArrayValue struct


.. type:: StructValue

   type StructValue map\[string\]Value


.. type:: NodeValue

   type NodeValue struct


.. type:: RelValue

   type RelValue struct


.. type:: RecursiveRelValue

   type RecursiveRelValue struct


.. type:: MapValue

   type MapValue struct


.. type:: UnionValue

   type UnionValue struct


.. type:: UUIDValue

   type UUIDValue uuid.UUID


.. type:: DecimalValue

   type DecimalValue decimal.Decimal


.. type:: ExternallyTagged

   type ExternallyTagged struct

   Json marshal helper for externally-tagged types  Mimics the externally taggged serde format in rust: https://serde.rs/enum-representations.html  examples: \`"Any"\` \`\{"Type": "content"\}\`


Interfaces
----------

.. type:: LogicalType

   **Methods:**

   * IsKuzuLogicalType()
   * MarshalJSON() (\[\]byte, error)


.. type:: Parameters

   **Methods:**

   * AsParameters() (map\[string\]Value, error)


.. type:: Value

   **Methods:**

   * IsKuzuValue()
   * MarshalJSON() (\[\]byte, error)


Variables
---------
