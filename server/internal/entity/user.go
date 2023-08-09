package entity

import "golang.org/x/crypto/bcrypt"

const (
	lenLogin    = 4
	lenPassword = 6
)

type User struct {
	ID       string
	Login    string
	Password string
}

func (u *User) IsValidPassword() bool {
	return u.Password != "" && len(u.Password) >= lenPassword
}

func (u *User) IsValidLogin() bool {
	return u.Password != "" && len(u.Login) >= lenLogin
}

func (u *User) SaveHashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) IsEqual(other User) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(other.Password))
	if err != nil {
		return false
	}

	return u.Login == other.Login
}
