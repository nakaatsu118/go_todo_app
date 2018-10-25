package handler

import (
	"fmt"

	"../model"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ShowFirstPage(c *gin.Context) {
	// fmt.Print("\x1b[2J") // 画面全クリア
	fmt.Println("-----------------------------")
	fmt.Println(" Welcome to Go-Todo service.")
	fmt.Println(" Please choose menu.")

	for {
		fmt.Println("-----------------------------")
		fmt.Println(" [Login    ]: Enter")
		fmt.Println(" [User Menu]: 'u' & Enter")
		fmt.Println(" [Exit     ]: 'e' & Enter")
		choice := model.StdIn()
		if choice == "" {
			ShowLogin(c)
		} else if choice == "u" {
			ShowUserMenu(c)
		} else if choice == "e" {
			fmt.Println("-----------------------------")
			fmt.Println(" Goodbyeﾉｼ")
			break
		} else {
			fmt.Println(" Please enter again.")
		}
	}

}

func ShowLogin(c *gin.Context) {
	fmt.Println("-----------------------------")
	fmt.Println(" Please login to your account.")
	fmt.Println("-----------------------------")

	fmt.Printf(" [Name]:")
	name := model.StdIn()
	fmt.Printf(" [Pass]:")
	pass := model.StdIn()
	encryptedPass := model.ToSha1Hash(pass)

	model.SetSessionLoginInfo(c, name, encryptedPass)
	Login(c)
}

func Login(c *gin.Context) {
	isExist, user := model.CheckLoginUserExist(c)
	if isExist == false {
		fmt.Println("-----------------------------")
		fmt.Println(" Failed to login.")
		fmt.Println(" Back to main menu.")
	} else {
		id := user.ID
		model.SetSessionUserInfo(c, id)
		fmt.Println("-----------------------------")
		fmt.Println(" Login successed!")
		ShowTodo(c)

	}
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	fmt.Println("-----------------------------")
	fmt.Println(" Logout finished.")
	fmt.Println(" Back to main menu.")
}

func ShowUserMenu(c *gin.Context) {
	fmt.Println("-----------------------------")
	fmt.Println(" User manage menu")
	fmt.Println(" Please choose menu.")

	for {
		fmt.Println("-----------------------------")
		fmt.Println(" [Add user   ]: 'a' & Enter")
		fmt.Println(" [Delete user]: 'd' & Enter")
		fmt.Println(" [Exit       ]: 'e' & Enter")
		choice := model.StdIn()
		if choice == "a" {
			AddUser(c)
		} else if choice == "d" {
			DeleteUser(c)
		} else if choice == "e" {
			fmt.Println("-----------------------------")
			fmt.Println(" Back to main menu.")
			break
		} else {
			fmt.Println("-----------------------------")
			fmt.Println(" Please enter again.")
		}
	}
}

func AddUser(c *gin.Context) {
	fmt.Printf(" Enter User name :")
	name := model.StdIn()
	fmt.Printf(" Enter Password  :")
	pass := model.StdIn()
	fmt.Printf(" Enter Password again:")
	passAgain := model.StdIn()
	if pass != passAgain {
		//fmt.Println(" Passwords do not match.")
		c.String(200, " Passwords do not match.")
	} else {
		model.AddUser(name, pass)
	}
}

func DeleteUser(c *gin.Context) {
	fmt.Printf(" Enter User name want to delete :")
	name := model.StdIn()
	fmt.Println(" Are you sure?")
	fmt.Println(" [OK    ]: 'ok' & Enter")
	fmt.Println(" [CANCEL]: Enter")
	confirm := model.StdIn()
	if confirm != "ok" {
		c.String(200, " Delete user canceled.")
	} else {
		model.DeleteUser(name)
	}

}
