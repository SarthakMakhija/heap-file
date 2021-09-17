package heap_file

type FreePageList struct {
	pageIds []uint32
}

func InitializeFreePageList(startingPageId uint32, pageCount int) *FreePageList {
	freePageList := &FreePageList{}
	pageId := startingPageId

	for index := 1; index <= pageCount; index++ {
		freePageList.pageIds = append(freePageList.pageIds, pageId)
		pageId = pageId + 1
	}
	return freePageList
}
