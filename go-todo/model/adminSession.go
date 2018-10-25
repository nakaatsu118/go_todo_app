package model

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type AdminSession struct {
	BaseModel
	AdminId      uint   `gorm:"type:int(20); not null"`
	OneTimeToken string `gorm:"type:varchar(100);unique;index;not null"`
}

func (AdminSession) TableName() string {
	return "admin_sessions"
}

func IsLogin(token string) bool {
	session := AdminSession{}

	db := ConnectGorm()
	defer db.Close()
	db.First(&session, "one_time_token = ?", token)
	if session.ID != 0 {
		return true
	} else {
		return false
	}
}

func (AdminSession) Create(adminId uint) (AdminSession, error) {
	session := AdminSession{}

	db := ConnectGorm()
	defer db.Close()

	session.AdminId = adminId
	session.OneTimeToken = RandString()
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()

	db.Create(&session)
	return session, db.Error
}
