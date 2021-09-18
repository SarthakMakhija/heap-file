package heap_file

import "github.com/SarthakMakhija/b-plus-tree/heap-file/page"

type SlottedPageCache struct {
	pageById map[uint32]*page.SlottedPage
}

func NewSlottedPageCache() *SlottedPageCache {
	return &SlottedPageCache{
		pageById: make(map[uint32]*page.SlottedPage),
	}
}

func (slottedPageCache *SlottedPageCache) Put(pageId uint32, page *page.SlottedPage) {
	slottedPageCache.pageById[pageId] = page
}

func (slottedPageCache *SlottedPageCache) GetBy(pageId uint32) (*page.SlottedPage, bool) {
	slottedPage, found := slottedPageCache.pageById[pageId]
	return slottedPage, found
}

func (slottedPageCache *SlottedPageCache) Evict(pageId uint32) {
	delete(slottedPageCache.pageById, pageId)
}
