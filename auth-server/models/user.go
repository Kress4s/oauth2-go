package models

type User struct {
	ID     int64  `gorm:"column:id"`
	Name   string `gorm:"column:id"`
	Avatar string `gorm:"column:id"`
}

func (u *User) TableName() string {
	return "user"
}

func GenerateTestUser() *User {
	return &User{
		ID:     100,
		Name:   "xys",
		Avatar: "image://",
	}
}
