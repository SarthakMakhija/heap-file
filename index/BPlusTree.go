package index

type BPlusTree struct {
	fileName      string
	pagePool      *PagePool
	pageHierarchy *PageHierarchy
}

const metaPageCount = 1
const rootPageCount = 1

func Create(options Options) (*BPlusTree, error) {
	indexFile, err := Open(options)
	pagePool := New(indexFile, options)

	if err != nil {
		return nil, err
	}
	tree := &BPlusTree{
		fileName:      options.FileName,
		pagePool:      pagePool,
		pageHierarchy: InstantiateHierarchy(pagePool),
	}
	if err := tree.create(options); err != nil {
		return nil, err
	}
	return tree, nil
}

func (tree BPlusTree) Get(key []byte) (KeyValuePair, bool, error) {
	return tree.pageHierarchy.Get(key)
}

func (tree *BPlusTree) create(options Options) error {
	if tree.pagePool.ContainsZeroPages() {
		return tree.initialize(options)
	}
	return nil
}

func (tree *BPlusTree) initialize(options Options) error {
	err := tree.pagePool.Allocate(metaPageCount + rootPageCount + options.PreAllocatedPagePoolSize)
	if err != nil {
		return err
	}
	return nil
}
