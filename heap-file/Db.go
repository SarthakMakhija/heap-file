package heap_file

import (
	"github.com/SarthakMakhija/b-plus-tree/heap-file/tuple"
	"os"
)

type Db struct {
	bufferPool   *BufferPool
	freePageList *FreePageList
	heapFile     *HeapFile
}

func Open(options HeapFileOptions) (*Db, error) {
	dbFile, err := openFile(options.FileName)
	if err != nil {
		return nil, err
	}

	bufferPool := NewBufferPool(dbFile, options)
	db := &Db{
		bufferPool: bufferPool,
	}
	if err := db.create(options); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *Db) Put(tuple *tuple.Tuple) (tuple.TupleId, error) {
	return db.heapFile.Put(tuple)
}

func (db *Db) GetBy(tupleId tuple.TupleId) *tuple.Tuple {
	return db.heapFile.GetBy(tupleId)
}

func openFile(fileName string) (*os.File, error) {
	fileMode := os.O_CREATE | os.O_RDWR
	return os.OpenFile(fileName, fileMode, 0644)
}

func (db *Db) create(options HeapFileOptions) error {
	if db.bufferPool.ContainsZeroPages() {
		return db.initialize(options)
	}
	return nil
}

func (db *Db) initialize(options HeapFileOptions) error {
	_, err := db.bufferPool.Allocate(options.PreAllocatedPagePoolSize)
	if err != nil {
		return err
	}
	db.freePageList = InitializeFreePageList(0, options.PreAllocatedPagePoolSize)
	db.heapFile = NewHeapFile(db.bufferPool, db.freePageList, options)
	return nil
}
