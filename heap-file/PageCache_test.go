package heap_file

import (
	"github.com/SarthakMakhija/heap-file/heap-file/page"
	"github.com/SarthakMakhija/heap-file/heap-file/tuple"
	"testing"
)

func TestPutsAndGetsASlottedPageInPageCache(t *testing.T) {
	slottedPageCache := NewSlottedPageCache()
	slottedPageCache.Put(100, page.NewSlottedPage(100, 100, tuple.TupleDescriptor{}))

	slottedPage, _ := slottedPageCache.GetBy(100)
	expectedPageId := uint32(100)

	if expectedPageId != slottedPage.PageId() {
		t.Fatalf("Expected cached page id to be %v, received %v", expectedPageId, slottedPage.PageId())
	}
}

func TestEvictsAPageFromPageCache(t *testing.T) {
	slottedPageCache := NewSlottedPageCache()
	slottedPageCache.Put(100, page.NewSlottedPage(100, 100, tuple.TupleDescriptor{}))
	slottedPageCache.Evict(100)

	_, found := slottedPageCache.GetBy(100)

	if found != false {
		t.Fatalf("Expected page to not be found in the page cache after eviction")
	}
}
