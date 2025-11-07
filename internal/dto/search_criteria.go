package dto

type Operator string

const (
	OpEqual     Operator = "="
	OpNotEqual  Operator = "!="
	OpGreater   Operator = ">"
	OpGreaterEq Operator = ">="
	OpLess      Operator = "<"
	OpLessEq    Operator = "<="
	OpIn        Operator = "in"
	OpLike      Operator = "like"
)

type SearchCriteria struct {
	Limit            int               `json:"limit" validate:"required,gte=0"`
	Offset           *int              `json:"offset" validate:"omitempty,gte=0"`
	OrderBy          *string           `json:"order_by"`
	SearchConditions []SearchCondition `json:"search_conditions"`
}

type SearchCondition struct {
	Field     string   `json:"field" validate:"required"`
	Operation Operator `json:"operation" validate:"required"`
	Value     any      `json:"value" validate:"required"`
}
