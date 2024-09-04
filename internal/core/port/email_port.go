package port

type EmailService interface {
	RegisterNotify(email string) error
	SetResetPasswordLink(email, encryptEmail string) error
	NewPasswordNotify(email, newRandomPassword string) error
}
