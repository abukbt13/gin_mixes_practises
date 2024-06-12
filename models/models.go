package models

type File struct {
	Id        uint     `form:"id" gorm:"primaryKey;autoIncrement"`
	Name      string   `form:"name"`
	Filenames []string `gorm:"type:json"`
}
