package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

type User struct {
	BaseModel
	Name     string `gorm:"type:varchar(50);not null;unique"`
	Password string `gorm:"type:varchar(200);not null;unique"`
}

type Users struct {
	Users []User `json:"users"`
}

type Todo struct {
	UserID int    `gorm:"type:int(20);not null;primary_key"`
	TodoID int    `gorm:"type:int(20);not null;primary_key"`
	Body   string `gorm:"type:varchar(200);not null"`
	// LimitDate interface{} `gorm:"column:limit_date" sql:"null;type:timestamp"`
	// LimitDate time.Time `gorm:"column:limit_date" sql:"null;type:timestamp"`
	LimitDate mysql.NullTime `gorm:"column:limit_date" sql:"null;type:timestamp"`
	CreatedAt time.Time      `gorm:"column:created_at" sql:"not null;type:timestamp"`
	UpdatedAt time.Time      `gorm:"column:updated_at" sql:"not null;type:timestamp"`
}

type Todos struct {
	Todos []Todo `json:"todos"`
}

func GetLoginUser(name interface{}, pass interface{}) (User, error) {
	user := User{}
	db := ConnectGorm()
	defer db.Close()
	if db.Find(&user, "name = ? and password = ?", name, pass); db.Error != nil {
		return user, db.Error
	} else {
		if db.RecordNotFound() {
			return user, errors.New("record not found")
		} else {
			return user, nil
		}
	}
}

func GetUserTodos(id interface{}) ([]Todo, error) {
	todos := []Todo{}
	db := ConnectGorm()
	defer db.Close()
	if db.Find(&todos, "user_id = ?", id); db.Error != nil {
		return todos, db.Error
	} else {
		if db.RecordNotFound() {
			return todos, nil
		} else {
			for k, m := range todos {
				todos[k] = m
			}
			return todos, nil
		}
	}
}

func AddTask(userID int, body string, limitDate interface{}, dateNull bool) (Todo, error) {
	task := Todo{}

	db := ConnectGorm()
	defer db.Close()

	task.UserID = userID
	task.Body = body

	if dateNull == true {
		task.LimitDate = mysql.NullTime{Time: time.Now(), Valid: false}
	} else {
		timeLimitDate := limitDate.(time.Time)
		task.LimitDate = mysql.NullTime{Time: timeLimitDate, Valid: true}
	}

	// fmt.Println("after:" + task.LimitDate.Time.String())

	db.Create(&task)

	fmt.Println("")
	fmt.Println("-----------------------------")
	fmt.Println(" Add task completed.")
	return task, db.Error
}

func DeleteTask(userID int, todoID int) (Todo, error) {
	task := Todo{}

	db := ConnectGorm()
	defer db.Close()

	if db.Model(&task).Where("user_id = ? and todo_id = ?", userID, todoID).Delete(&task); db.Error != nil {
		return task, errors.New("record not found")
	} else {
		fmt.Println("")
		fmt.Println("-----------------------------")
		fmt.Println(" Delete task completed.")
		return task, nil
	}
}

func AddUser(name string, pass string) (User, error) {
	user := User{}

	db := ConnectGorm()
	defer db.Close()

	user.Name = name
	user.Password = ToSha1Hash(pass)
	db.Create(&user)

	return user, db.Error
}

func DeleteUser(name string) (User, error) {
	user := User{}

	db := ConnectGorm()
	defer db.Close()

	if db.First(&user, "name=?", name).RecordNotFound() {
		fmt.Println(" record not found")
		return user, errors.New("record not found")
	} else {
		if db.Model(&user).Where("name = ?", name).Delete(&user); db.Error != nil {
			return user, db.Error
		} else {
			fmt.Println("")
			fmt.Println("-----------------------------")
			fmt.Println(" Delete user completed.")
			return user, nil
		}
	}
}
