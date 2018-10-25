package model

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	SessionToken = "SessionToken"
	SessionKey   = "session"
)

type BaseModel struct {
	ID        int       `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"column:created_at" sql:"not null;type:timestamp"`
	UpdatedAt time.Time `gorm:"column:updated_at" sql:"not null;type:timestamp"`
}

type DatabaseInfo struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}

func ConnectGorm() *gorm.DB {
	var db DatabaseInfo
	// db.Host = "localhost"
	// db.Username = "root"
	// db.Password = ""
	// db.Port = "3306"
	// db.Database = "go_todo"
	db.Host = "mysql"
	db.Username = "root"
	db.Password = "root"
	db.Port = "3306"
	db.Database = "go_todo"

	dbCon := fmt.Sprintf("%s:%s@tcp([%s]:%s)/%s?charset=utf8&parseTime=True",
		db.Username,
		db.Password,
		db.Host,
		db.Port,
		db.Database)

	// dbCon := fmt.Sprintf("%s:%s@tcp([%s]:%s)/%s?charset=utf8&parseTime=True",
	// 	"root",
	// 	"",
	// 	"127.0.0.1",
	// 	"3306",
	// 	"go_todo")

	if db, err := gorm.Open("mysql", dbCon); err != nil {
		panic("failed to connect database")
	} else {
		return db
	}
}

func ToSha1Hash(before string) string {
	hash := sha1.New()
	hash.Write([]byte(before))
	after := hash.Sum(nil)
	return fmt.Sprintf("%x", after)
}

func StdIn() string {
	stdin := bufio.NewScanner(os.Stdin)
	stdin.Scan()
	text := stdin.Text()
	return text
}
