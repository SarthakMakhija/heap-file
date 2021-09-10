package index

type GetResult struct {
	keyValuePair KeyValuePair
	index        int
	page         *Page
	found        bool
	err          error
}

func NewKeyAvailableGetResult(pair KeyValuePair, index int, page *Page) GetResult {
	return GetResult{
		keyValuePair: pair,
		index:        index,
		page:         page,
		found:        true,
		err:          nil,
	}
}

func NewKeyMissingGetResult(index int, page *Page) GetResult {
	return GetResult{
		keyValuePair: KeyValuePair{},
		index:        index,
		page:         page,
		found:        false,
		err:          nil,
	}
}

func NewFailedGetResult(err error) GetResult {
	return GetResult{
		found: false,
		err:   err,
	}
}
