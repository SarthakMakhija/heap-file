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

func (pageHierarchy PageHierarchy) Get(key []byte) GetResult {
	return pageHierarchy.get(key, pageHierarchy.rootPage)
}

func (pageHierarchy PageHierarchy) RootPageId() int {
	return pageHierarchy.rootPage.id
}

func (pageHierarchy PageHierarchy) PageById(id int) *Page {
	return pageHierarchy.pageById[id]
}

func (pageHierarchy PageHierarchy) get(key []byte, page *Page) GetResult {
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
		child, err := pageHierarchy.fetchPage(page.childPageIds[index])
		if err != nil {
			return NewFailedGetResult(err)
		}
		return pageHierarchy.get(key, child)
	}
}

func (pageHierarchy PageHierarchy) fetchPage(pageId int) (*Page, error) {
	page, found := pageHierarchy.pageById[pageId]
	if found {
		return page, nil
	}
	page, err := pageHierarchy.pagePool.Read(pageId)
	if err != nil {
		return nil, err
	}
	pageHierarchy.pageById[page.id] = page
	return page, nil
}
