package dto

// Operator represents comparison operators for search conditions
// @Name Operator
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

// SearchCriteria represents search criteria with pagination and filtering
// @Name SearchCriteria
type SearchCriteria struct {
	Limit            int               `json:"limit" validate:"required,gte=0"`
	Offset           *int              `json:"offset" validate:"omitempty,gte=0"`
	OrderBy          *string           `json:"order_by" validate:"omitempty,min=1,max=50"`
	SearchConditions []SearchCondition `json:"search_conditions" validate:"omitempty,dive"`
}

// SearchCondition represents a single search condition
// @Name SearchCondition
type SearchCondition struct {
	Field     string   `json:"field" validate:"required"`
	Operation Operator `json:"operation" validate:"required"`
	Value     any      `json:"value" validate:"required"`
}
