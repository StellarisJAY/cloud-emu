package common

import "gorm.io/gorm"

type Pagination struct {
	Page     int32
	PageSize int32
	Total    int32
}

func WithPagination(p *Pagination) func(db *gorm.DB) *gorm.DB {
	limit := int(p.PageSize)
	offset := int((p.Page - 1) * p.PageSize)
	return func(db *gorm.DB) *gorm.DB {
		var total int64
		db.Count(&total)
		p.Total = int32(total)
		return db.Limit(limit).Offset(offset)
	}
}
