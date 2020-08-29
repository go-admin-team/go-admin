package dto

type Dtor interface {
	Validate() error
	Generate() Dtor
	GetPageIndex() int
	GetPageSize() int
	GetNeedSearch() interface{}
}
