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

func (tree BPlusTree) Get(key []byte) (KeyValuePair, bool, error) {
	return tree.get(key, tree.rootPage)
}

func (tree BPlusTree) get(key []byte, page *Page) (KeyValuePair, bool, error) {
	index, found := page.Get(key)
	if page.isLeaf() {
		if found {
			return page.GetKeyValuePairAt(index), found, nil
		}
		return KeyValuePair{}, false, nil
	}
	//handle if found in non-leaf??
	child, err := tree.fetchPage(page.childPageIds[index])
	if err != nil {
		return KeyValuePair{}, false, err
	}
	return tree.get(key, child)
}

func (tree BPlusTree) fetchPage(pageId int) (*Page, error) {
	page, found := tree.pageById[pageId]
	if found {
		return page, nil
	}
	page, err := tree.pagePool.Read(pageId)
	if err != nil {
		return nil, err
	}
	tree.pageById[page.id] = page
	return page, nil
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
