package heap_file

import "os"

type Db struct {
	bufferPool   *BufferPool
	freePageList *FreePageList
}

func Open(options Options) (*Db, error) {
	dbFile, err := openFile(options.FileName)
	if err != nil {
		return nil, err
	}

	bufferPool := NewBufferPool(dbFile, options.PageSize)
	db := &Db{
		bufferPool: bufferPool,
	}
	if err := db.create(options); err != nil {
		return nil, err
	}
	return db, nil
}

func openFile(fileName string) (*os.File, error) {
	fileMode := os.O_CREATE | os.O_RDWR
	return os.OpenFile(fileName, fileMode, 0644)
}

func (db *Db) create(options Options) error {
	if db.bufferPool.ContainsZeroPages() {
		return db.initialize(options)
	}
	return nil
}

func (db *Db) initialize(options Options) error {
	_, err := db.bufferPool.Allocate(options.PreAllocatedPagePoolSize)
	if err != nil {
		return err
	}
	db.freePageList = InitializeFreePageList(0, options.PreAllocatedPagePoolSize)
	return nil
}
