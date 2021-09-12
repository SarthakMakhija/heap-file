package index

import (
	"bytes"
)

type PageHierarchy struct {
	rootPage                       *Page
	pageById                       map[int]*Page
	pagePool                       *PagePool
	allowedPageOccupancyPercentage int
	freePageList                   *FreePageList
}

func NewPageHierarchy(pagePool *PagePool, allowedPageOccupancyPercentage int, freePageList *FreePageList) *PageHierarchy {
	pageHierarchy := &PageHierarchy{
		rootPage:                       NewPage(0),
		pagePool:                       pagePool,
		pageById:                       map[int]*Page{},
		allowedPageOccupancyPercentage: allowedPageOccupancyPercentage,
		freePageList:                   freePageList,
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

	if pageHierarchy.isPageEligibleForSplit(pageHierarchy.rootPage) {
		if err := splitRoot(); err != nil {
			return err
		}
	}
	if err := pageHierarchy.put(keyValuePair, pageHierarchy.rootPage); err != nil {
		return err
	}
	pageHierarchy.Write()
	return nil
}

func (pageHierarchy *PageHierarchy) Get(key []byte) GetResult {
	return pageHierarchy.get(key, pageHierarchy.rootPage)
}

func (pageHierarchy *PageHierarchy) Write() {
	for _, page := range pageHierarchy.pageById {
		if page.IsDirty() {
			pageHierarchy.pagePool.Write(page)
		}
	}
}

func (pageHierarchy PageHierarchy) RootPageId() int {
	return pageHierarchy.rootPage.id
}

func (pageHierarchy PageHierarchy) PageById(id int) *Page {
	return pageHierarchy.pageById[id]
}

func (pageHierarchy *PageHierarchy) put(keyValuePair KeyValuePair, page *Page) error {
	if page.isLeaf() {
		index, _ := page.Get(keyValuePair.key)
		//assume not found
		page.insertAt(index, keyValuePair)
		return nil
	}
	return pageHierarchy.insertOrSplit(keyValuePair, page)
}

func (pageHierarchy *PageHierarchy) insertOrSplit(keyValuePair KeyValuePair, page *Page) error {
	index, found := page.Get(keyValuePair.key)
	if found {
		index = index + 1
	}

	childPage, err := pageHierarchy.fetchOrCachePage(page.childPageIds[index])
	if err != nil {
		return err
	}
	if pageHierarchy.isPageEligibleForSplit(childPage) {
		sibling, err := pageHierarchy.allocateSinglePage()
		if err != nil {
			return err
		}
		if err := childPage.split(page, sibling, index); err != nil {
			return err
		}
		if bytes.Compare(keyValuePair.key, page.keyValuePairs[index].key) >= 0 {
			childPage, err = pageHierarchy.fetchOrCachePage(page.childPageIds[index+1])
			if err != nil {
				return err
			}
		}
	}
	return pageHierarchy.put(keyValuePair, childPage)
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

func (pageHierarchy PageHierarchy) isPageEligibleForSplit(page *Page) bool {
	return page.size() >= (pageHierarchy.allowedPageOccupancyPercentage * (pageHierarchy.pagePool.pageSize) / 100)
}

func (pageHierarchy *PageHierarchy) allocateSinglePage() (*Page, error) {
	pages, err := pageHierarchy.allocatePages(1)
	if err != nil {
		return nil, err
	}
	return pages[0], nil
}

func (pageHierarchy *PageHierarchy) allocatePages(pageCount int) ([]*Page, error) {
	newPageId := pageHierarchy.freePageList.allocateAndUpdate(pageCount)
	if newPageId < 1 {
		var err error
		newPageId, err = pageHierarchy.pagePool.Allocate(pageCount)
		if err != nil {
			return nil, err
		}
	}
	pages := make([]*Page, pageCount)
	for index := 0; index < pageCount; index++ {
		newPage := NewPage(newPageId)
		pageHierarchy.pageById[newPageId] = newPage
		pages[index] = newPage
		newPageId = newPageId + 1
	}
	return pages, nil
}
