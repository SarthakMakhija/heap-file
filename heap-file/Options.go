package heap_file

import (
	"github.com/SarthakMakhija/b-plus-tree/heap-file/field"
	"github.com/SarthakMakhija/b-plus-tree/heap-file/tuple"
	"os"
)

type Options struct {
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

func DefaultOptions() Options {
	return Options{
		PageSize:                 os.Getpagesize(),
		FileName:                 "heap.db",
		PreAllocatedPagePoolSize: 10,
		TupleDescriptor: tuple.TupleDescriptor{
			FieldTypes: []field.FieldType{field.StringFieldType{}, field.Uint16FieldType{}},
		},
	}
}
