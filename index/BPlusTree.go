package index

type BPlusTree struct {
	fileName      string
	pagePool      *PagePool
	pageHierarchy *PageHierarchy
}

const metaPageCount = 1
const rootPageCount = 1

func CreateBPlusTree(options Options) (*BPlusTree, error) {
	indexFile, err := OpenIndexFile(options)
	pagePool := NewPagePool(indexFile, options)

	if err != nil {
		return nil, err
	}
	tree := &BPlusTree{
		fileName:      options.FileName,
		pagePool:      pagePool,
		pageHierarchy: NewPageHierarchy(pagePool),
	}
	if err := tree.create(options); err != nil {
		return nil, err
	}
	return tree, nil
}

func (tree BPlusTree) Put(key, value []byte) {
	tree.pageHierarchy.Put(KeyValuePair{key: key, value: value})
}

func (tree BPlusTree) Get(key []byte) GetResult {
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
