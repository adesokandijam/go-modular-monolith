package shared

import "strings"

type ProductFilters struct {
	Sort     string `validate:"oneof=name -name price_in_pence -price_in_pence created_at -created_at stock_quantity -stock_quantity"`
	PageSize int64  `validate:"gte=1,lte=100"`
	Page     int64  `validate:"gte=1"`
}

// func ValidateFilters(v *validator.Validate, pf ProductFilters){

// }

func (pf ProductFilters) SortColumns() string {
	return strings.TrimPrefix(pf.Sort, "-")
}

func (pf ProductFilters) SortDirection() string {
	if strings.HasPrefix(pf.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (pf ProductFilters) Limit() int64 {
	return pf.PageSize
}

func (pf ProductFilters) Offset() int64 {
	return (pf.Page - 1) * pf.PageSize
}
