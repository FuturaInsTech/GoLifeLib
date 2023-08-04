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

type ReportData struct {
	Title              string
	ColumnHeadings     []string
	DataRows           [][]interface{}
	CommonDataHeads    []string
	CommonDataDetails  []string
	FooterDataHeads    []string
	FooterDataDetails  []string
	TemplateName       string
	AdditionalData     interface{}
	FirstPageRecCount  int
	SubseqPageRecCount int
	ReportType         string
	FileName           string
}

func (r ReportData) ConvertTo2D(d1array []string, count int) [][]string {

	d2array := make([][]string, 0)

	for i := 0; i < len(d1array); {
		tempArr := make([]string, 0)

		for j := 0; j < count; j++ {

			if i >= len(d1array) {
				break
			}
			tempArr = append(tempArr, d1array[i])
			i++

		}
		d2array = append(d2array, tempArr)
	}

	return d2array
}

func (r ReportData) Get1Dindex(index1 int, index2 int, count int) int {

	return ((index1 * count) + index2)
}

func (r ReportData) GetDetailPages() [][][]interface{} {

	d3array := make([][][]interface{}, 0)
	pageNum := 1
	for i := 0; i < len(r.DataRows); {
		tempArr := make([][]interface{}, 0)
		recordsPerPage := r.SubseqPageRecCount
		if pageNum == 1 {
			recordsPerPage = r.FirstPageRecCount
		}

		for j := 0; j < recordsPerPage; j++ {

			if i >= len(r.DataRows) {
				break
			}
			tempArr = append(tempArr, r.DataRows[i])
			i++

		}
		d3array = append(d3array, tempArr)
		pageNum++
	}

	return d3array

}
