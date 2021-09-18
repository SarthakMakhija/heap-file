package tuple

type MarshalledTuple struct {
	buffer []byte
	size   int
}

func (marshalledTuple MarshalledTuple) Buffer() []byte {
	return marshalledTuple.buffer
}

func (marshalledTuple MarshalledTuple) Size() int {
	return marshalledTuple.size
}
