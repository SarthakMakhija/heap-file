package index

import "os"

type IndexFile struct {
	file      *os.File
	pageSize  int
	pageCount int
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
	indexFile.pageCount = indexFile.numberOfPages()
	return indexFile, nil
}

func (indexFile *IndexFile) numberOfPages() int {
	fileSize, _ := indexFile.fileSize()
	return int(fileSize) / indexFile.pageSize
}

func (indexFile *IndexFile) fileSize() (int64, error) {
	stat, err := indexFile.file.Stat()
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}
