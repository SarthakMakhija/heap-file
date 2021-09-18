package heap_file

import "os"

type Options struct {
	// PageSize for file I/O. All reads and writes will always
	// be done with pages of this size. Must be multiple of os.Getpagesize().
	PageSize int

	// Name of the db file
	FileName string

	// PreAllocatedPagePoolSize identifies the number of pages to be pre-allocated when the Db is opened.
	// Must be greater than 0
	PreAllocatedPagePoolSize int
}

func DefaultOptions() Options {
	return Options{
		PageSize:                 os.Getpagesize(),
		FileName:                 "heap.db",
		PreAllocatedPagePoolSize: 10,
	}
}
