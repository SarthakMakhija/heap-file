package index

type PageHierarchy struct {
	rootPage *Page
	pageById map[int]*Page
	pagePool *PagePool
}

func NewPageHierarchy(pagePool *PagePool) *PageHierarchy {
	pageHierarchy := &PageHierarchy{
		rootPage: NewPage(0),
		pagePool: pagePool,
		pageById: map[int]*Page{},
	}
	pageHierarchy.pageById[pageHierarchy.rootPage.id] = pageHierarchy.rootPage
	return pageHierarchy
}

func (pageHierarchy *PageHierarchy) Put(keyValuePair KeyValuePair) error {
	splitRoot := func() error {
		siblingPageCount := 1
		newRootPageCount := 1

		pages, err := pageHierarchy.allocatePages(siblingPageCount + newRootPageCount)
		if err != nil {
			return err
		}
		newRootPage, rightSiblingPage, oldRootPage := pages[0], pages[1], pageHierarchy.rootPage
		newRootPage.childPageIds = append(newRootPage.childPageIds, oldRootPage.id)
		pageHierarchy.rootPage = newRootPage

		return oldRootPage.split(newRootPage, rightSiblingPage, 0)
	}

	if len(pageHierarchy.rootPage.keyValuePairs) >= 3 { //will be replaced with % occupancy later
		if err := splitRoot(); err != nil {
			return err
		}
	}
	getResult := pageHierarchy.Get(keyValuePair.key)
	//assume key does not exist and page is not full
	getResult.page.insertAt(getResult.index, keyValuePair)
	return nil
}

func (pageHierarchy *PageHierarchy) Get(key []byte) GetResult {
	return pageHierarchy.get(key, pageHierarchy.rootPage)
}

func (pageHierarchy PageHierarchy) RootPageId() int {
	return pageHierarchy.rootPage.id
}

func (pageHierarchy PageHierarchy) PageById(id int) *Page {
	return pageHierarchy.pageById[id]
}

func (pageHierarchy *PageHierarchy) get(key []byte, page *Page) GetResult {
	index, found := page.Get(key)
	if page.isLeaf() {
		if found {
			return NewKeyAvailableGetResult(page.GetKeyValuePairAt(index), index, page)
		}
		return NewKeyMissingGetResult(index, page)
	} else {
		if found {
			index = index + 1
		}
		child, err := pageHierarchy.fetchOrCachePage(page.childPageIds[index])
		if err != nil {
			return NewFailedGetResult(err)
		}
		return pageHierarchy.get(key, child)
	}
}

func (pageHierarchy *PageHierarchy) fetchOrCachePage(pageId int) (*Page, error) {
	page, found := pageHierarchy.pageById[pageId]
	if found {
		return page, nil
	}
	page, err := pageHierarchy.pagePool.Read(pageId)
	if err != nil {
		return nil, err
	}
	pageHierarchy.pageById[pageId] = page
	return page, nil
}

func (pageHierarchy *PageHierarchy) allocatePages(pageCount int) ([]*Page, error) {
	newPageId := 1 //will come from free list of pages
	pages := make([]*Page, pageCount)
	for index := 0; index < pageCount; index++ {
		newPage := NewPage(newPageId)
		pageHierarchy.pageById[newPageId] = newPage
		pages[index] = newPage
		newPageId = newPageId + 1
	}
	return pages, nil
}
