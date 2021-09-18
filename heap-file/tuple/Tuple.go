package tuple

import "github.com/SarthakMakhija/heap-file/heap-file/field"

type Tuple struct {
	fields []field.Field
}

func NewTuple() *Tuple {
	var fields []field.Field
	return &Tuple{
		fields: fields,
	}
}

func (tuple *Tuple) AddField(field field.Field) {
	tuple.fields = append(tuple.fields, field)
}

func (tuple Tuple) MarshalBinary() MarshalledTuple {
	var buffer []byte
	for _, aField := range tuple.fields {
		buffer = append(buffer, aField.MarshalBinary()...)
	}
	return MarshalledTuple{buffer: buffer, size: len(buffer)}
}

func (tuple Tuple) Size() int {
	return tuple.MarshalBinary().Size()
}

func (tuple *Tuple) UnMarshalBinary(buffer []byte, fieldTypes []field.FieldType) {
	offset := 0
	for _, fieldType := range fieldTypes {
		aField := fieldType.UnMarshalBinary(buffer[offset:])
		offset = offset + aField.MarshalSize()
		tuple.fields = append(tuple.fields, aField)
	}
}

func (tuple Tuple) AllFields() []field.Field {
	return tuple.fields
}

func (tuple Tuple) KeyField() field.Field {
	//for now it just returns the last field as a key field
	return tuple.fields[len(tuple.fields)-1]
}
