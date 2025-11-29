package purchase

// AddToCartRequest represents request for adding product to cart
// @Name AddToCartRequest
type AddToCartRequest struct {
	ProductID uint `json:"product_id" validate:"required,min=1"`
	Quantity  uint `json:"quantity" validate:"required,min=1"`
}

// PurchaseRequest represents request for purchasing cart
// @Name PurchaseRequest
type PurchaseRequest struct {
	CartID uint `json:"cart_id" validate:"required,min=1"`
}
