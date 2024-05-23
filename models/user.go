package models

type User struct {
	Id       uint   `form:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `form:"name"`
	Email    string `form:"email" gorm:"unique"`
	Password string `form:"password"`
}
