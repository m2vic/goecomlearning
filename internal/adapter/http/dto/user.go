package dto

type UpdateUserRequest struct {
	EmailRequest
	FirstName      string `json:"firstname" validate:"required"`
	LastName       string `json:"lastname" validate:"required"`
	AddressDetails string `json:"address"`
}
type LoginRequest struct {
	Username string `json:"username" validate:"min=5"`
	Password string `json:"password" validate:"min=5"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required"`
	Username string `json:"username" validate:"required,min=5"`
	Password string `json:"password" validate:"required,min=5"`
}
type ChangePasswordRequest struct {
	OldPassword string `json:"oldpassword" validate:"min=5"`
	NewPassword string `json:"newpassword" validate:"min=5"`
}
type EmailRequest struct {
	Email string `json:"email" validate:"required"`
}

type CheckoutRequest struct {
	Product []Product `json:"productlist" validate:"required"`
}

type DeleteItemInCartRequest struct {
	Product
}
type IncreaseItemInCartRequest struct {
	Product
}
type DecreaseItemInCartRequest struct {
	Product
}
