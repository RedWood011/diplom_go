package entity

type Secret struct {
	ID     string
	UserID string
	Data   []byte
	Name   []byte
}
