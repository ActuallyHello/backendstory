package common

import (
	"log/slog"

	"github.com/ActuallyHello/backendstory/internal/dto"
	"gorm.io/gorm"
)

func BuildQuery(queryCtx *gorm.DB, criteria dto.SearchCriteria) *gorm.DB {
	return queryCtx.
		Scopes(applySearchConditions(criteria.SearchConditions)).
		Scopes(applyPagination(criteria.Limit, criteria.Offset)).
		Scopes(applyOrdering(criteria.OrderBy))
}

func applySearchConditions(conditions []dto.SearchCondition) func(db *gorm.DB) *gorm.DB {
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

func applyCondition(queryCtx *gorm.DB, condition dto.SearchCondition) *gorm.DB {
	switch condition.Operation {
	case dto.OpEqual:
		return queryCtx.Where(condition.Field+" = ?", condition.Value)
	case dto.OpNotEqual:
		return queryCtx.Where(condition.Field+" != ?", condition.Value)
	case dto.OpGreater:
		return queryCtx.Where(condition.Field+" > ?", condition.Value)
	case dto.OpGreaterEq:
		return queryCtx.Where(condition.Field+" >= ?", condition.Value)
	case dto.OpLess:
		return queryCtx.Where(condition.Field+" < ?", condition.Value)
	case dto.OpLessEq:
		return queryCtx.Where(condition.Field+" <= ?", condition.Value)
	case dto.OpIn:
		return queryCtx.Where(condition.Field+" in ?", condition.Value)
	case dto.OpLike:
		return queryCtx.Where(condition.Field+" like ?", condition.Value)
	default:
		slog.Warn("Undefined operation for search criteria", "operation", condition.Operation)
		return queryCtx
	}
}
