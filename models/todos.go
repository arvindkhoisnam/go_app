package models

import "gorm.io/gorm"

type Todos struct {
	ID 		   uint    `gorm:"primary key;autoIncrement" json:"id"`
	Completed  *bool   `json:"completed"`
	Body 	   *string `json:"body"`
}

func MigrateTodos(db *gorm.DB) error {
	err := db.AutoMigrate(&Todos{})
	return err
}