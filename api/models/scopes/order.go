package scopes

import (
	"fmt"
	"goa-golang/helpers"
	"goa-golang/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormOrderer interface {
	ToOrder(tableName string, defaultField string, vars []interface{}, orderByOptions ...string) func(db *gorm.DB) *gorm.DB
}

type GormOrder struct {
	*utils.Order
}

func (r *GormOrder) ToOrder(tableName string, defaultField string, vars []interface{}, orderByOptions ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !helpers.InArray(r.Order.GetOrderBy(), orderByOptions) {
			return db.Clauses(clause.OrderBy{
				Expression: clause.Expr{SQL: fmt.Sprintf("FIELD(`%v`.`%v`,?) %v", tableName, defaultField, r.Order.GetSortBy()), Vars: vars, WithoutParentheses: true},
			})
		}
		return db.Order(fmt.Sprintf("`%v`.`%v` %v", tableName, r.Order.GetOrderBy(), r.Order.GetSortBy()))
	}
}
