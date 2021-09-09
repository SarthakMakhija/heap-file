package index

import (
	"github.com/edsrzf/mmap-go"
	"os"
)

type IndexFile struct {
	file         *os.File
	size         int64
	pageSize     int
	pageCount    int
	memoryMapped mmap.MMap
}

func Open(options Options) (*IndexFile, error) {
	fileMode := os.O_CREATE | os.O_RDWR
	file, err := os.OpenFile(options.FileName, fileMode, 0644)

	if err != nil {
		return nil, err
	}
	indexFile := &IndexFile{
		file:     file,
		pageSize: options.PageSize,
	}
	indexFile.size, _ = indexFile.fileSize()
	indexFile.pageCount = indexFile.numberOfPages()
	return indexFile, nil
}

func (indexFile IndexFile) ContainsZeroPages() bool {
	return indexFile.pageCount == 0
}

func (indexFile *IndexFile) Allocate(pageCount int) error {
	err := indexFile.unMap()
	if err != nil {
		return err
	}
	targetSize := indexFile.size + int64(pageCount*indexFile.pageSize)
	if err := indexFile.file.Truncate(targetSize); err != nil {
		return err
	}

	indexFile.size = targetSize
	indexFile.pageCount = indexFile.numberOfPages()

	return indexFile.mMap()
}

func (indexFile *IndexFile) numberOfPages() int {
	return int(indexFile.size) / indexFile.pageSize
}

func (indexFile *IndexFile) fileSize() (int64, error) {
	stat, err := indexFile.file.Stat()
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}

func (indexFile *IndexFile) unMap() error {
	if indexFile.file == nil || indexFile.memoryMapped == nil {
		return nil
	}
	return indexFile.memoryMapped.Unmap()
}

func (indexFile *IndexFile) mMap() error {
	if err := indexFile.unMap(); err != nil {
		return err
	}
	memoryMapped, err := mmap.Map(indexFile.file, mmap.RDWR, 0)
	if err != nil {
		return err
	}
	indexFile.memoryMapped = memoryMapped
	return nil
}
