package errs

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

var (
	DatabaseNotFound  = &Error{Code: 400, Message: "missing database"}
	UserAlreadyExist  = &Error{Code: 400, Message: "user already exists"}
	UsernameInvalid   = &Error{Code: 400, Message: "username invalid"}
	PasswordInvalid   = &Error{Code: 400, Message: "password invalid"}
	EmailAlreadyExist = &Error{Code: 400, Message: "email already exist"}
	EmailNotFound     = &Error{Code: 400, Message: "email not found"}
	UpdateUserFail    = &Error{Code: 500, Message: "fail to update user"}
	GenerateTokenFail = &Error{Code: 500, Message: "fail to generate tokens"}
	HashPasswordFail  = &Error{Code: 500, Message: "fail to hash password"}
	TokenNotFound     = &Error{Code: 400, Message: "token not found"}

	//product
	ProductNotFound  = &Error{Code: 400, Message: "product not found"}
	NotEnoughProduct = &Error{Code: 400, Message: "insufficient products in stock"}
)
