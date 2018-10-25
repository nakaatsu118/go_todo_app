package model

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SetSessionLoginInfo(c *gin.Context, name string, pass string) {
	session := sessions.Default(c)
	session.Set("name", name)
	session.Set("password", pass)
	session.Save()
}

func SetSessionUserInfo(c *gin.Context, id int) {
	// fmt.Println("userID:" + strconv.Itoa(id))
	session := sessions.Default(c)
	// fmt.Println("name:" + session.Get("name").(string))
	/* ここセッション保持できるかが変わる */
	// session.Set("id", strconv.Itoa(id))
	session.Set("id", id)
	// fmt.Println("id:" + session.Get("id").(string))
	session.Save()
	// fmt.Println("session id:" + session.Get("id").(string))
}

func CheckLoginUserExist(c *gin.Context) (bool, User) {
	session := sessions.Default(c)
	name := session.Get("name")
	pass := session.Get("password")

	if user, err := GetLoginUser(name, pass); err != nil {
		return false, user
	} else {
		return user.ID != 0, user
	}
}
