package common

import (
	"gorm.io/gorm"
)

// SearchOperator тип для операторов сравнения
type SearchOperator string

const (
	OpEqual     SearchOperator = "="
	OpNotEqual  SearchOperator = "!="
	OpLike      SearchOperator = "LIKE"
	OpILike     SearchOperator = "ILIKE" // case-insensitive LIKE
	OpIn        SearchOperator = "IN"
	OpNotIn     SearchOperator = "NOT IN"
	OpGt        SearchOperator = ">"
	OpGte       SearchOperator = ">="
	OpLt        SearchOperator = "<"
	OpLte       SearchOperator = "<="
	OpIsNull    SearchOperator = "IS NULL"
	OpIsNotNull SearchOperator = "IS NOT NULL"
)

// SearchCondition условие для поиска
type SearchCondition struct {
	Field     string         `json:"field"`
	Operator  SearchOperator `json:"operator"`
	Value     interface{}    `json:"value,omitempty"`
	ValueList []interface{}  `json:"valueList,omitempty"` // для IN/NOT IN
}

// OrderDirection направление сортировки
type OrderDirection string

const (
	OrderASC  OrderDirection = "ASC"
	OrderDESC OrderDirection = "DESC"
)

// OrderBy параметры сортировки
type OrderBy struct {
	Field     string         `json:"field"`
	Direction OrderDirection `json:"direction"`
}

// SearchCriteria критерии поиска
type SearchCriteria struct {
	Conditions []SearchCondition `json:"conditions,omitempty"`
	OrderBy    []OrderBy         `json:"orderBy,omitempty"`
	Limit      int               `json:"limit,omitempty"`
	Offset     int               `json:"offset,omitempty"`
	WithTotal  bool              `json:"withTotal,omitempty"` // возвращать ли общее количество
	Distinct   bool              `json:"distinct,omitempty"`  // использовать DISTINCT
}

// Apply применяет критерии поиска к GORM запросу
func (sc *SearchCriteria) Apply(query *gorm.DB) *gorm.DB {
	// Применяем условия
	for _, condition := range sc.Conditions {
		query = sc.applyCondition(query, condition)
	}

	// Применяем сортировку
	for _, order := range sc.OrderBy {
		query = query.Order(order.Field + " " + string(order.Direction))
	}

	// Применяем пагинацию
	if sc.Limit > 0 {
		query = query.Limit(sc.Limit)
	}
	if sc.Offset > 0 {
		query = query.Offset(sc.Offset)
	}

	// Применяем DISTINCT
	if sc.Distinct {
		query = query.Distinct()
	}

	return query
}

// applyCondition применяет одно условие к запросу
func (sc *SearchCriteria) applyCondition(query *gorm.DB, condition SearchCondition) *gorm.DB {
	switch condition.Operator {
	case OpEqual, OpNotEqual, OpGt, OpGte, OpLt, OpLte:
		return query.Where(condition.Field+" "+string(condition.Operator)+" ?", condition.Value)
	case OpLike, OpILike:
		value := condition.Value
		if strVal, ok := value.(string); ok {
			value = "%" + strVal + "%"
		}
		if condition.Operator == OpILike {
			// Для MySQL ILIKE не поддерживается, используем LOWER
			return query.Where("LOWER("+condition.Field+") LIKE LOWER(?)", value)
		}
		return query.Where(condition.Field+" LIKE ?", value)
	case OpIn:
		return query.Where(condition.Field+" IN (?)", condition.ValueList)
	case OpNotIn:
		return query.Where(condition.Field+" NOT IN (?)", condition.ValueList)
	case OpIsNull:
		return query.Where(condition.Field + " IS NULL")
	case OpIsNotNull:
		return query.Where(condition.Field + " IS NOT NULL")
	default:
		return query
	}
}

// AddCondition добавляет условие поиска
func (sc *SearchCriteria) AddCondition(field string, operator SearchOperator, value interface{}) *SearchCriteria {
	sc.Conditions = append(sc.Conditions, SearchCondition{
		Field:    field,
		Operator: operator,
		Value:    value,
	})
	return sc
}

// AddOrder добавляет сортировку
func (sc *SearchCriteria) AddOrder(field string, direction OrderDirection) *SearchCriteria {
	sc.OrderBy = append(sc.OrderBy, OrderBy{
		Field:     field,
		Direction: direction,
	})
	return sc
}

// SetPagination устанавливает пагинацию
func (sc *SearchCriteria) SetPagination(limit, offset int) *SearchCriteria {
	sc.Limit = limit
	sc.Offset = offset
	return sc
}

// SetWithTotal включает возврат общего количества
func (sc *SearchCriteria) SetWithTotal(withTotal bool) *SearchCriteria {
	sc.WithTotal = withTotal
	return sc
}
