package utils

// PageInfo is the pagination response, matching Java PageInfo<T>
type PageInfo[T any] struct {
	Total    int64  `json:"total"`
	List     []T    `json:"list"`
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
	Pages    int    `json:"pages"`
}

func NewPageInfo[T any](total int64, list []T, pageNum, pageSize int) *PageInfo[T] {
	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}
	return &PageInfo[T]{
		Total:    total,
		List:     list,
		PageNum:  pageNum,
		PageSize: pageSize,
		Pages:    pages,
	}
}
