package scopes

import (
	"goa-golang/helpers"
	"goa-golang/utils"

	"gorm.io/gorm"
)

type GormPager interface {
	ToPaginate() func(db *gorm.DB) *gorm.DB
}

type GormPagination struct {
	*utils.Pagination
}

func (r *GormPagination) ToPaginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(helpers.OffsetCal(r.Pagination.GetPage(), r.Pagination.GetLimit())).Limit(r.Pagination.GetLimit())
	}
}
