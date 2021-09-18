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

func (freePageList *FreePageList) allocateAndUpdate(pages int) (bool, uint32) {
	isAvailable, firstFreePageId, remainingFreePageIds := freePageList.allocateContiguous(pages)
	freePageList.pageIds = remainingFreePageIds
	return isAvailable, firstFreePageId
}

func (freePageList *FreePageList) allocateContiguous(pages int) (bool, uint32, []uint32) {
	if len(freePageList.pageIds) < pages {
		return false, 0, freePageList.pageIds
	} else if pages == 1 {
		return true, freePageList.pageIds[0], freePageList.pageIds[1:]
	}

	startingIndex, endIndex := 0, 0
	for ; startingIndex < len(freePageList.pageIds); startingIndex++ {
		endIndex = startingIndex + (pages - 1)
		if endIndex < len(freePageList.pageIds) && freePageList.pageIds[endIndex] == freePageList.pageIds[startingIndex]+uint32(pages-1) {
			break
		}
	}

	if startingIndex >= len(freePageList.pageIds) || endIndex >= len(freePageList.pageIds) {
		return false, 0, freePageList.pageIds
	}

	firstFreePageId := freePageList.pageIds[startingIndex]
	freePageList.pageIds = append(freePageList.pageIds[:startingIndex], freePageList.pageIds[endIndex+1:]...)
	return true, firstFreePageId, freePageList.pageIds
}
