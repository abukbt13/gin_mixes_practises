package models

type Picture struct {
	Id          uint   `form:"id" gorm:"primaryKey;autoIncrement"`
	Description string `form:"description"`
	PictureName string `form:"picture_name"`
}
