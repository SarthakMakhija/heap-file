package index

type FreePageList struct {
	pageIds []int
}

func InitializeFreePageList(startingPageId int, pageCount int) *FreePageList {
	freePageList := &FreePageList{}
	pageId := startingPageId

	for index := 1; index <= pageCount; index++ {
		freePageList.pageIds = append(freePageList.pageIds, pageId)
		pageId = pageId + 1
	}
	return freePageList
}
