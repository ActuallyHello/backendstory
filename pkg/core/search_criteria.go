package core

import (
	"log/slog"
	"strings"

	"gorm.io/gorm"
)

func BuildQuery(queryCtx *gorm.DB, criteria SearchCriteria) *gorm.DB {
	return queryCtx.
		Scopes(applySearchConditions(criteria.SearchConditions)).
		Scopes(applyPagination(criteria.Limit, criteria.Offset)).
		Scopes(applyOrdering(criteria.OrderBy))
}

func applySearchConditions(conditions []SearchCondition) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, condition := range conditions {
			db = applyCondition(db, condition)
		}
		return db
	}
}

func applyPagination(limit int, offset *int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Limit(limit)
		if offset != nil {
			db = db.Offset(*offset)
		}
		return db
	}
}

func applyOrdering(orderBy *string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if orderBy != nil && *orderBy != "" {
			db = db.Order(*orderBy)
		}
		return db
	}
}

func applyCondition(queryCtx *gorm.DB, condition SearchCondition) *gorm.DB {
	field := convertToDBField(condition.Field)
	switch condition.Operation {
	case OpEqual:
		return queryCtx.Where(field+" = ?", condition.Value)
	case OpNotEqual:
		return queryCtx.Where(field+" != ?", condition.Value)
	case OpGreater:
		return queryCtx.Where(field+" > ?", condition.Value)
	case OpGreaterEq:
		return queryCtx.Where(field+" >= ?", condition.Value)
	case OpLess:
		return queryCtx.Where(field+" < ?", condition.Value)
	case OpLessEq:
		return queryCtx.Where(field+" <= ?", condition.Value)
	case OpIn:
		return queryCtx.Where(field+" in ?", condition.Value)
	case OpLike:
		return queryCtx.Where(field+" like ?", condition.Value)
	default:
		slog.Warn("Undefined operation for search criteria", "operation", condition.Operation)
		return queryCtx
	}
}

func convertToDBField(input string) string {
	input = strings.ReplaceAll(input, "_", "")
	return strings.ToUpper(input)
}
