package main

import (
	"./handler"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	// session管理
	store := sessions.NewCookieStore([]byte("secret"))
	server.Use(sessions.Sessions("todoSessions", store))

	// 初期画面
	server.GET("/", handler.ShowFirstPage)

	// ログイン・ログアウト
	server.POST("/checkUser", handler.CheckUser)
	server.GET("/login", handler.ShowLogin)
	server.POST("/login", handler.Login)
	server.GET("/logout", handler.Logout)

	// ユーザー管理
	server.GET("/user", handler.ShowUserMenu)
	server.POST("/adduser", handler.AppAddUser)
	server.POST("/deleteuser", handler.AppDeleteUser)

	// Todo表示
	server.GET("/todo", handler.ShowTodo)
	server.POST("/addtask", handler.AppAddTask)
	server.POST("/deletetask", handler.AppDeleteTask)

	// server.GET("/checkTodo", handler.CheckTodo)
	server.POST("/checkTodo", handler.CheckTodo)

	// server.GET("/countTodo", handler.CountTodo)
	server.POST("/countTodo", handler.CountTodo)

	server.POST("/getTodoBody", handler.GetTodoBody)
	server.POST("/getTodoLimitDate", handler.GetTodoLimitDate)

	server.POST("/returnUID", handler.ReturnUID)

	// ポート番号指定
	server.Run(":1192")
}
