package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"../model"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginUser struct {
	UserID   int    `form:"id"  json:"id" xml:"id" binding:"exists"`
	Name     string `form:"name" json:"name" xml:"name" binding:"exists"`
	Password string `form:"password" json:"password" xml:"password" binding:"exists"`
}

type Todo struct {
	UserID int    `form:"id"  json:"id" xml:"id" binding:"exists"`
	Num    string `form:"num" json:"num" xml:"num" binding:"exists"`
}

type AddTodo struct {
	UserID    int    `form:"id"  json:"id" xml:"id" binding:"exists"`
	Body      string `form:"task" json:"task" xml:"task" binding:"exists"`
	LimitDate string `form:"limitDate" json:"limitDate" xml:"limitDate" binding:"exists"`
}

func CheckUser(c *gin.Context) {
	// リクエストボディを型にバインディング
	req := LoginUser{}

	if err := c.Bind(&req); err != nil {
		log.Fatal(err)

		// c.String(200, "err")
	} else {
		name := req.Name
		pass := req.Password

		session := sessions.Default(c)
		//session := sessions.GetSessionsInfo(c)
		session.Set("name", name)
		session.Set("password", pass)
		session.Save()

		result, user := model.CheckLoginUserExist(c)

		if result == false {
			c.String(http.StatusOK, "ng")
		} else {
			c.String(http.StatusOK, "ok")
			model.SetSessionUserInfo(c, user.ID)
			// fmt.Println("session id:" + session.Get("id").(string))
		}

	}
}

func CheckTodo(c *gin.Context) {
	req := LoginUser{}

	if err := c.Bind(&req); err != nil {
		log.Fatal(err)

		// c.String(200, "err")
	} else {
		userID := req.UserID

		// fmt.Println("userID:" + strconv.Itoa(userID))
		// session := sessions.Default(c)
		// userID := session.Get("id")
		// fmt.Println("name:" + session.Get("name").(string))

		// fmt.Println("user ID:" + userID.(string))

		if _, err := model.GetUserTodos(userID); err != nil {
			c.String(http.StatusOK, "ng")
		} else {
			c.String(http.StatusOK, "ok")
		}
	}
}

func CountTodo(c *gin.Context) {
	req := LoginUser{}

	if err := c.Bind(&req); err != nil {
		log.Fatal(err)

		// c.String(200, "err")
	} else {
		userID := req.UserID

		// session := sessions.Default(c)
		// userID := session.Get("id")

		if length, err := model.CountTodo(userID); err != nil {
			c.String(http.StatusOK, "ng")
		} else {
			if length == 0 {
				c.String(http.StatusOK, "zero")
			} else {
				lengthStr := strconv.Itoa(length)
				c.String(http.StatusOK, lengthStr)
			}
		}
	}
}

func GetTodoBody(c *gin.Context) {
	// リクエストボディを型にバインディング
	req := Todo{}

	if err := c.Bind(&req); err != nil {
		log.Fatal(err)
		// c.String(200, "err")
	} else {
		numStr := req.Num
		id := req.UserID
		// fmt.Println("numStr:" + numStr)
		num, _ := strconv.Atoi(numStr)
		// session := sessions.Default(c)
		// userID := session.Get("id")

		// fmt.Println("userID:" + id.(string))

		if todos, err := model.GetUserTodos(id); err != nil {
			c.String(http.StatusOK, "ng")
		} else {
			body := todos[num].Body
			c.String(http.StatusOK, body)
		}
	}

}

func GetTodoLimitDate(c *gin.Context) {
	// リクエストボディを型にバインディング
	req := Todo{}

	if err := c.Bind(&req); err != nil {
		log.Fatal(err)
		// c.String(200, "err")
	} else {
		numStr := req.Num
		id := req.UserID
		num, _ := strconv.Atoi(numStr)
		// session := sessions.Default(c)
		// userID := session.Get("id")

		if todos, err := model.GetUserTodos(id); err != nil {
			c.String(http.StatusOK, "ng")
		} else {
			limitDate := todos[num].LimitDate.Time.String()
			// limitDate := todos[num].LimitDate.(string)
			// limitDate := todos[num].LimitDate.String()
			// fmt.Println("limitData:" + limitDate)
			c.String(http.StatusOK, limitDate)
		}
	}
}

func ReturnUID(c *gin.Context) {
	// リクエストボディを型にバインディング
	req := LoginUser{}

	if err := c.Bind(&req); err != nil {
		log.Fatal(err)

		// c.String(200, "err")
	} else {
		name := req.Name
		pass := req.Password

		session := sessions.Default(c)
		//session := sessions.GetSessionsInfo(c)
		session.Set("name", name)
		session.Set("password", pass)
		session.Save()

		result, user := model.CheckLoginUserExist(c)

		if result == false {
			c.String(http.StatusOK, "ng")
		} else {
			uID := strconv.Itoa(user.ID)
			c.String(http.StatusOK, uID)
			// model.SetSessionUserInfo(c, user.ID)
			// fmt.Println("session id:" + session.Get("id").(string))
		}
	}

}

func AppAddTask(c *gin.Context) {
	// リクエストボディを型にバインディング
	var limitDate interface{}
	var dateNull bool
	layout := "2006-01-02 15:04:05"
	req := AddTodo{}

	if err := c.Bind(&req); err != nil {
		log.Fatal(err)
		// c.String(200, "err")
	} else {
		id := req.UserID
		task := req.Body
		limitDateString := req.LimitDate
		// fmt.Println("before:" + limitDateString)
		if limitDateString == "0001-01-01 00:00:00" {
			// limitDate = mysql.NullTime{Time: time.Now(), Valid: false}
			limitDate = time.Time{}
			dateNull = true
		} else {
			limitDate, _ = time.Parse(layout, limitDateString)
			dateNull = false
		}

		if todos, err := model.AddTask(id, task, limitDate, dateNull); err != nil {
			c.String(http.StatusOK, "ng")
		} else {
			if todos.Body != "" {
				c.String(http.StatusOK, "ok")
			}
		}
	}
}

func AppDeleteTask(c *gin.Context) {
	// リクエストボディを型にバインディング
	req := Todo{}

	if err := c.Bind(&req); err != nil {
		log.Fatal(err)
		// c.String(200, "err")
	} else {
		deleteIDStr := req.Num
		userID := req.UserID
		deleteID, _ := strconv.Atoi(deleteIDStr)

		if todos, err := model.GetUserTodos(userID); err != nil {
			c.String(http.StatusOK, "get todos failed")
		} else {
			// limitDate := todos[num].LimitDate.Time.String()
			todoID := todos[deleteID-1].TodoID
			if todos, err := model.DeleteTask(userID, todoID); err != nil {
				c.String(http.StatusOK, "ng")
			} else {
				if todos.UserID != 0 {
					c.String(http.StatusOK, "ok")
				}
			}
		}
	}
}

func AppAddUser(c *gin.Context) {
	req := LoginUser{}

	if err := c.Bind(&req); err != nil {
		log.Fatal(err)
	} else {
		name := req.Name
		pass := req.Password

		if users, err := model.AddUser(name, pass); err != nil {
			c.String(http.StatusOK, "ng")
		} else {
			if users.ID != 0 {
				c.String(http.StatusOK, "ok")
			}
		}
	}
}

func AppDeleteUser(c *gin.Context) {
	req := LoginUser{}

	if err := c.Bind(&req); err != nil {
		log.Fatal(err)
	} else {
		name := req.Name

		if users, err := model.DeleteUser(name); err != nil {
			c.String(http.StatusOK, "ng")
		} else {
			if users.ID != 0 {
				c.String(http.StatusOK, "ok")
			}
		}
	}
}
