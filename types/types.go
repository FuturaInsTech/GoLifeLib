package types

type User struct {
	Id          uint64
	Name        string
	Email       string
	Phone       string
	UserGroupId uint
	LanguageId  uint
	CompanyId   uint
}

type PageFilterSort struct {
	Offset      int
	PageSize    int
	Filters     []Filter
	SortingList string
}

type Sorting struct {
	Id   string
	Desc bool
}

type Filter struct {
	Id    string
	Value string
}
