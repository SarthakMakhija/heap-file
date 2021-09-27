package heap_file

import (
	"os"
)

type OverflowPagePool struct {
	file     *os.File
	pageSize int
}

func NewOverflowPool(file *os.File, pageSize int) *OverflowPagePool {
	return &OverflowPagePool{
		file:     file,
		pageSize: pageSize,
	}
}

func (overflowPagePool OverflowPagePool) Read(pageId uint32) (*OverflowPage, error) {
	readPage := func(pageId uint32) (*OverflowPage, error) {
		buffer := make([]byte, overflowPagePool.pageSize)
		_, err := overflowPagePool.file.ReadAt(buffer, overflowPagePool.offsetOf(pageId))
		if err != nil {
			return nil, err
		}
		return NewReadonlyOverflowPage(buffer), nil
	}

	return readPage(pageId)
}

func (overflowPagePool *OverflowPagePool) Write(page *OverflowPage) error {
	_, err := overflowPagePool.file.WriteAt(page.buffer, overflowPagePool.offsetOf(page.id))
	if err != nil {
		return err
	}
	return nil
}

func (overflowPagePool OverflowPagePool) offsetOf(pageId uint32) int64 {
	return int64(uint32(overflowPagePool.pageSize) * pageId)
}

func (overflowPagePool OverflowPagePool) numberOfPages() int {
	return overflowPagePool.fileSize() / overflowPagePool.pageSize
}

func (overflowPagePool OverflowPagePool) fileSize() int {
	stat, err := overflowPagePool.file.Stat()
	if err != nil {
		return 0 //Handle later
	}
	return int(stat.Size())
}
