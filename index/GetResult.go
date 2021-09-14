package index

type GetResult struct {
	keyValuePair KeyValuePair
	index        int
	page         *Page
	found        bool
	Err          error
}

func NewKeyAvailableGetResult(pair KeyValuePair, index int, page *Page) GetResult {
	return GetResult{
		keyValuePair: pair,
		index:        index,
		page:         page,
		found:        true,
		Err:          nil,
	}
}

func NewKeyMissingGetResult(index int, page *Page) GetResult {
	return GetResult{
		keyValuePair: KeyValuePair{},
		index:        index,
		page:         page,
		found:        false,
		Err:          nil,
	}
}

func NewFailedGetResult(err error) GetResult {
	return GetResult{
		found: false,
		Err:   err,
	}
}
