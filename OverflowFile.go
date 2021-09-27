package heap_file

type OverflowFile struct {
	overflowPagePool    *OverflowPagePool
	currentOverflowPage *OverflowPage
}

func NewOverflowFile(overflowPagePool *OverflowPagePool, pageSize int) *OverflowFile {
	return &OverflowFile{
		overflowPagePool:    overflowPagePool,
		currentOverflowPage: NewOverflowPage(0, pageSize),
	}
}

func (overflowFile *OverflowFile) Put(buffer []byte) *Slot {
	slot := overflowFile.currentOverflowPage.Put(buffer)
	overflowFile.writeAllPages()
	return slot
}

func (overflowFile *OverflowFile) GetAt(slotNo int) []byte {
	return overflowFile.getAt(overflowFile.currentOverflowPage, slotNo, []byte{})
}

func (overflowFile *OverflowFile) getAt(overflowPage *OverflowPage, slotNo int, buffer []byte) []byte {
	content, nextOverflowSlotId := overflowPage.GetAt(slotNo)
	if nextOverflowSlotId != 0 {
		nextOverflowPageId := overflowPage.NextOverflowPageId()
		nextOverflowPage, err := overflowFile.overflowPagePool.Read(nextOverflowPageId)
		if err != nil {
			panic(err)
		}
		return overflowFile.getAt(nextOverflowPage, int(nextOverflowSlotId), append(buffer, content...))
	}
	return append(buffer, content...)
}

func (overflowFile *OverflowFile) writeAllPages() {
	page := overflowFile.currentOverflowPage
	for page != nil {
		err := overflowFile.overflowPagePool.Write(page)
		if err != nil {
			panic(err)
		}
		page = page.nextPage
	}
}
