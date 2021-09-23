package schema

import (
	"io"
	"time"
	"unsafe"
)

var (
	_ = unsafe.Sizeof(0)
	_ = io.ReadFull
	_ = time.Now()
)

type PersistentTupleId struct {
	PageId uint32
	SlotNo uint32
}

func (d *PersistentTupleId) Size() (s uint64) {

	s += 8
	return
}
func (d *PersistentTupleId) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{

		buf[0+0] = byte(d.PageId >> 0)

		buf[1+0] = byte(d.PageId >> 8)

		buf[2+0] = byte(d.PageId >> 16)

		buf[3+0] = byte(d.PageId >> 24)

	}
	{

		buf[0+4] = byte(d.SlotNo >> 0)

		buf[1+4] = byte(d.SlotNo >> 8)

		buf[2+4] = byte(d.SlotNo >> 16)

		buf[3+4] = byte(d.SlotNo >> 24)

	}
	return buf[:i+8], nil
}

func (d *PersistentTupleId) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{

		d.PageId = 0 | (uint32(buf[0+0]) << 0) | (uint32(buf[1+0]) << 8) | (uint32(buf[2+0]) << 16) | (uint32(buf[3+0]) << 24)

	}
	{

		d.SlotNo = 0 | (uint32(buf[0+4]) << 0) | (uint32(buf[1+4]) << 8) | (uint32(buf[2+4]) << 16) | (uint32(buf[3+4]) << 24)

	}
	return i + 8, nil
}
