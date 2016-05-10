package proto

type Category struct {
	Name    string           `json:"name,omitempty"`
	Details *CategoryDetails `json:"details,omitempty"`
}

type CategoryDetails struct {
	DetailsCreation
}

func NewCategoryRequest() Request {
	return Request{
		Category: &Category{},
	}
}

func NewCategoryResult() Result {
	return Result{
		Errors:     &[]string{},
		Categories: &[]Category{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
