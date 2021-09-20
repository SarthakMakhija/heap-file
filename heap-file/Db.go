package heap_file

import (
	"github.com/SarthakMakhija/heap-file/heap-file/field"
	"github.com/SarthakMakhija/heap-file/heap-file/tuple"
	"github.com/SarthakMakhija/heap-file/index"
	"os"
)

type Db struct {
	bufferPool   *BufferPool
	freePageList *FreePageList
	heapFile     *HeapFile
	bPlusTree    *index.BPlusTree
}

func Open(options DbOptions) (*Db, error) {
	dbFile, err := openFile(options.FileName())
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
	tupleId, err := db.heapFile.Put(tuple)
	if err == nil {
		err = db.bPlusTree.Put(tuple.KeyField().MarshalBinary(), tupleId.MarshalBinary())
		return tupleId, err
	}
	return tupleId, err
}

func (db *Db) GetByKey(key field.Field) (*tuple.Tuple, error) {
	getResult := db.bPlusTree.Get(key.MarshalBinary())
	if getResult.Err == nil {
		if len(getResult.KeyValuePair.RawValue()) != 0 {
			tupleId := &tuple.TupleId{}
			tupleId.UnMarshalBinary(getResult.KeyValuePair.RawValue())
			return db.GetByTupleId(*tupleId), nil
		}
		return tuple.NewTuple(), nil
	}
	return nil, getResult.Err
}

func (db *Db) GetByTupleId(tupleId tuple.TupleId) *tuple.Tuple {
	return db.heapFile.GetBy(tupleId)
}

func openFile(fileName string) (*os.File, error) {
	fileMode := os.O_CREATE | os.O_RDWR
	return os.OpenFile(fileName, fileMode, 0644)
}

func (db *Db) create(options DbOptions) error {
	if db.bufferPool.ContainsZeroPages() {
		return db.initialize(options)
	}
	return nil
}

func (db *Db) initialize(options DbOptions) error {
	_, err := db.bufferPool.Allocate(options.PreAllocatedPagePoolSize())
	if err != nil {
		return err
	}
	bPlusTree, err := index.CreateBPlusTree(options.IndexOptions)
	if err != nil {
		return err
	}

	db.freePageList = InitializeFreePageList(0, options.PreAllocatedPagePoolSize())
	db.heapFile = NewHeapFile(db.bufferPool, db.freePageList, options)
	db.bPlusTree = bPlusTree
	return nil
}
