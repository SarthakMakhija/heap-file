package index

type BPlusTree struct {
	fileName string
	pagePool *PagePool
	rootPage *Page
	pageById map[int]*Page
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
		fileName: options.FileName,
		pagePool: pagePool,
		rootPage: nil,
		pageById: map[int]*Page{},
	}
	if err := tree.create(options); err != nil {
		return nil, err
	}
	return tree, nil
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
	tree.rootPage = NewPage(0)
	tree.pageById[tree.rootPage.id] = tree.rootPage
	return nil
}
