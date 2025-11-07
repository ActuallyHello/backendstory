package common

import (
	"log/slog"

	"github.com/ActuallyHello/backendstory/internal/dto"
	"gorm.io/gorm"
)

func BuildQuery(queryCtx *gorm.DB, criteria dto.SearchCriteria) *gorm.DB {
	if criteria.SearchConditions != nil {
		for _, condition := range criteria.SearchConditions {
			queryCtx = applyCondition(queryCtx, condition)
		}
	}

	queryCtx.Limit(criteria.Limit)
	if criteria.Offset != nil {
		queryCtx.Offset(*criteria.Offset)
	}

	if criteria.OrderBy != nil {
		queryCtx.Order(criteria.OrderBy)
	}

	return queryCtx
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
