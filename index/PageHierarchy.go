package index

type PageHierarchy struct {
	rootPage *Page
	pageById map[int]*Page
	pagePool *PagePool
}

func InstantiateHierarchy(pagePool *PagePool) *PageHierarchy {
	pageHierarchy := &PageHierarchy{
		rootPage: NewPage(0),
		pagePool: pagePool,
		pageById: map[int]*Page{},
	}
	pageHierarchy.pageById[pageHierarchy.rootPage.id] = pageHierarchy.rootPage
	return pageHierarchy
}

func (pageHierarchy PageHierarchy) Get(key []byte) (KeyValuePair, bool, error) {
	var getInner func(*Page) (KeyValuePair, bool, error)
	getInner = func(page *Page) (KeyValuePair, bool, error) {
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
			return getInner(child)
		}
	}
	return getInner(pageHierarchy.rootPage)
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
