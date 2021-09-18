package heap_file

import (
	"github.com/SarthakMakhija/b-plus-tree/heap-file/field"
	"github.com/SarthakMakhija/b-plus-tree/heap-file/tuple"
	"github.com/SarthakMakhija/b-plus-tree/index"
	"os"
)

type DbOptions struct {
	HeapFileOptions HeapFileOptions
	IndexOptions    index.Options
}

func (dbOptions DbOptions) PageSize() int {
	return dbOptions.HeapFileOptions.PageSize
}

func (dbOptions DbOptions) PreAllocatedPagePoolSize() int {
	return dbOptions.HeapFileOptions.PreAllocatedPagePoolSize
}

func (dbOptions DbOptions) TupleDescriptor() tuple.TupleDescriptor {
	return dbOptions.HeapFileOptions.TupleDescriptor
}

func (dbOptions DbOptions) FileName() string {
	return dbOptions.HeapFileOptions.FileName
}

type HeapFileOptions struct {
	// PageSize for file I/O. All reads and writes will always
	// be done with pages of this size. Must be multiple of os.Getpagesize().
	PageSize int

	// Name of the db file
	FileName string

	// PreAllocatedPagePoolSize identifies the number of pages to be pre-allocated when the Db is opened.
	// Must be greater than 0
	PreAllocatedPagePoolSize int

	//For now, TupleDescriptor is a part of options, which later (in the actual implementation) will move to
	//an abstraction say Catalog
	TupleDescriptor tuple.TupleDescriptor
}

func DefaultOptions() DbOptions {
	return DbOptions{
		HeapFileOptions: HeapFileOptions{
			PageSize:                 os.Getpagesize(),
			FileName:                 "heap.db",
			PreAllocatedPagePoolSize: 10,
			TupleDescriptor: tuple.TupleDescriptor{
				FieldTypes: []field.FieldType{field.StringFieldType{}, field.Uint16FieldType{}},
			},
		},
		IndexOptions: index.DefaultOptions(),
	}
}
