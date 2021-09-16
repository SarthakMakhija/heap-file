package field

type Field interface {
	MarshalBinary() []byte
}
