package index

type Page struct {
	id int
}

func NewPage(id int) *Page {
	return &Page{
		id: id,
	}
}
