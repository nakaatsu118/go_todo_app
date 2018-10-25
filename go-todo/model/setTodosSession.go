package model

import (
	"strconv"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SetSessionTodos(c *gin.Context, todos []Todo) {
	session := sessions.Default(c)
	length := len(todos)
	for i := 0; i < length; i++ {
		num := strconv.Itoa(i)
		name := "todos" + num
		session.Set(name+"todoID", todos[i].TodoID)
		session.Set(name+"body", todos[i].Body)
		session.Set(name+"limitDate", todos[i].LimitDate)
	}

	session.Save()
}

func CountTodo(id interface{}) (int, error) {
	todos := []Todo{}
	db := ConnectGorm()
	defer db.Close()
	if db.Find(&todos, "user_id = ?", id); db.Error != nil {
		return 0, db.Error
	} else {
		if db.RecordNotFound() {
			return 0, nil
		} else {
			length := len(todos)
			return length, nil
		}
	}
}
