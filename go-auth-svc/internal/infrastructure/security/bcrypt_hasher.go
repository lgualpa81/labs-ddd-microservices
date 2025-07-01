package security

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct{}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (h *BcryptHasher) Hash(plainText string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	return string(bytes), err
}

func (h *BcryptHasher) Compare(hashedText string, plainText string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedText), []byte(plainText))
	return err == nil
}
