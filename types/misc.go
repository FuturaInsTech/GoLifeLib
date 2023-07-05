package types

type SearchPagination struct {
	SearchString   string
	SearchCriteria string
	SortColumn     string
	SortDirection  string `gorm:"default:'asc'";`
	FirstTime      bool
	Offset         int
	PageNum        int
	PageSize       int
}
