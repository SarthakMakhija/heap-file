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

func (pageHierarchy PageHierarchy) Get(key []byte) (KeyValuePair, bool, error) {
	return pageHierarchy.get(key, pageHierarchy.rootPage)
}

func (pageHierarchy PageHierarchy) get(key []byte, page *Page) (KeyValuePair, bool, error) {
	index, found := page.Get(key)
	if page.isLeaf() {
		if found {
			return page.GetKeyValuePairAt(index), found, nil
		}
		return KeyValuePair{}, false, nil
	} else {
		if found {
			index = index + 1
		}
		child, err := pageHierarchy.fetchPage(page.childPageIds[index])
		if err != nil {
			return KeyValuePair{}, false, err
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
