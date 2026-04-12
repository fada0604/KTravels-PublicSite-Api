package pagination

type PageInfo struct {
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
	HasNext  bool  `json:"hasNext"`
	HasPrev  bool  `json:"hasPrev"`
}

func New(page, pageSize int, total int64) *PageInfo {
	offset := int64(page * pageSize)
	return &PageInfo{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		HasNext:  offset+int64(pageSize) < total,
		HasPrev:  page > 0,
	}
}
