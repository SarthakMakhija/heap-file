package tuple

import (
	"github.com/SarthakMakhija/heap-file/heap-file/schema"
)

type TupleId struct {
	PageId uint32
	SlotNo int //serialized as uint32
}

func (tupleId TupleId) MarshalBinary() []byte {
	persistentTupleId := &schema.PersistentTupleId{
		PageId: tupleId.PageId,
		SlotNo: uint32(tupleId.SlotNo),
	}
	buffer, _ := persistentTupleId.Marshal(nil)
	return buffer
}

func (tupleId *TupleId) UnMarshalBinary(buffer []byte) {
	persistentTupleId := &schema.PersistentTupleId{}
	_, _ = persistentTupleId.Unmarshal(buffer)

	tupleId.PageId = persistentTupleId.PageId
	tupleId.SlotNo = int(persistentTupleId.SlotNo)
}
