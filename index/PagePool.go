package index

type PagePool struct {
	indexFile *IndexFile
	pageSize  int
	pageCount int
}

func New(indexFile *IndexFile, options Options) *PagePool {
	pagePool := &PagePool{
		indexFile: indexFile,
	}
	pagePool.pageSize = options.PageSize
	pagePool.pageCount = pagePool.numberOfPages()

	return pagePool
}

func (pagePool *PagePool) Allocate(pageCount int) error {
	targetSize := pagePool.indexFile.size + int64(pageCount*pagePool.pageSize)
	if err := pagePool.indexFile.ResizeTo(targetSize); err != nil {
		return err
	}
	pagePool.pageCount = pagePool.numberOfPages()
	return nil
}

func (pagePool PagePool) Read(pageId int) ([]byte, error) {
	offset := func() int64 {
		return int64(pagePool.pageSize * pageId)
	}
	return pagePool.indexFile.readFrom(offset(), pagePool.pageSize)
}

func (pagePool PagePool) numberOfPages() int {
	return int(pagePool.indexFile.size) / pagePool.pageSize
}

func (pagePool PagePool) ContainsZeroPages() bool {
	return pagePool.pageCount == 0
}
