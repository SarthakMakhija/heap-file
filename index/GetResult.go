package index

type GetResult struct {
	keyValuePair KeyValuePair
	index        int
	pageId       int
	found        bool
	err          error
}

func NewKeyAvailableGetResult(pair KeyValuePair, index int, pageId int) GetResult {
	return GetResult{
		keyValuePair: pair,
		index:        index,
		pageId:       pageId,
		found:        true,
		err:          nil,
	}
}

func NewKeyMissingGetResult(index int, pageId int) GetResult {
	return GetResult{
		keyValuePair: KeyValuePair{},
		index:        index,
		pageId:       pageId,
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
