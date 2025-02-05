package models

import (
	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

type Todos struct {
	ID 		   uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Completed  *bool   `json:"completed"`
	Body 	   *string `json:"body"`
}

type User struct {
	ID 		  string 	`gorm:"primaryKey" json:"id"`
	Username  string	`json:"username"`
	Password  string	`json:"password"`

}
func (u *User)BeforeCreate(tx *gorm.DB) (err error) {
    u.ID = cuid.New()
    return
}

func MigrateDB(db *gorm.DB) error {
	err := db.AutoMigrate(&Todos{},&User{})
	return err
}
