package handler

import (
	"fmt"
	"strconv"
	"time"

	"../model"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type UserTodo struct {
	ID        int
	Body      string
	LimitDate time.Time
}

func ShowTodo(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("id")

	// fmt.Println("userID:" + userID.(string))
	for {
		if todos, err := model.GetUserTodos(userID); err != nil {
			c.String(200, " Failed to get todos.")
		} else {
			fmt.Println(" Here is your task.")
			length := len(todos)
			for i := 0; i < length; i++ {
				var limitDate time.Time
				fmt.Println("-----------------------------")
				fmt.Println("ID: " + strconv.Itoa(i+1))
				fmt.Println("Task: " + todos[i].Body)
				// if todos[i].LimitDate != nil {
				// limitDate = todos[i].LimitDate.(time.Time)
				limitDate = todos[i].LimitDate.Time
				// limitDate = todos[i].LimitDate
				if limitDate.IsZero() != true {
					fmt.Println("Limit: " + limitDate.Format("Mon Jan _2 2006"))
				}
				// }
			}
			fmt.Println("-----------------------------")
			fmt.Println("")
			fmt.Println(" [Add task    ]: Enter")
			fmt.Println(" [Delete task ]: 'd' & Enter")
			fmt.Println(" [Logout      ]: anykey & Enter")
			inputKey := model.StdIn()
			if inputKey == "d" {
				DeleteTask(c, userID, todos)
			} else if inputKey == "" {
				AddTask(c, userID)
			} else {
				Logout(c)
				break
			}
		}
	}

}

func AddTask(c *gin.Context, userID interface{}) {
	var limitDate interface{}
	layout := "2006-01-02 15:04:05"

	fmt.Printf("Task: ")
	body := model.StdIn()

	fmt.Println(" [Set limitdate        ]: Enter")
	fmt.Println(" [Don't set limitdate  ]: anykey & Enter")
	fmt.Println("")
	addTask := model.StdIn()
	if addTask == "" {
		fmt.Printf(" Limit Year(ex. 2018年 → 2018): ")
		limitYear := model.StdIn()
		fmt.Printf(" Limit Month(ex. 1月 → 01): ")
		limitMonth := model.StdIn()
		fmt.Printf(" Limit Day(ex. 2日 → 02): ")
		limitDay := model.StdIn()
		limitDateString := limitYear + "-" + limitMonth + "-" + limitDay + " 00:00:00"
		limitDate, _ = time.Parse(layout, limitDateString)
	} else {
		limitDate = mysql.NullTime{Time: time.Now(), Valid: false}
	}

	userIDInt := userID.(int)

	// model.AddTask(userIDInt, body, limitDate)
	model.AddTask(userIDInt, body, limitDate, true)

}

func DeleteTask(c *gin.Context, userID interface{}, todos []model.Todo) {
	fmt.Printf(" [Enter task ID to delete ]:")
	deleteIDString := model.StdIn()
	deleteID, _ := strconv.Atoi(deleteIDString)

	userIDInt := userID.(int)
	length := len(todos)
	var todoID int
	if deleteID > length {
		fmt.Println(" Enter ID is wrong.")
	} else if deleteID <= 0 {
		fmt.Println(" Enter ID is wrong.")
	} else {
		todoID = todos[deleteID-1].TodoID
		model.DeleteTask(userIDInt, todoID)
	}

}
