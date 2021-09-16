package field

type Field interface {
	MarshalBinary() []byte
	Value() interface{}
}
