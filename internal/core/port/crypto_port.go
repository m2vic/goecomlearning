package port

type CryptoService interface {
	Encrypt(email string) (string, error)
	Decrypt(email string) (string, error)
}
