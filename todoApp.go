package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	// urlBase = "http://localhost:1192"
	urlBase = "http://localhost:9999"
)

var directory string

type Todo struct {
	UserID    int         `gorm:"type:int(20);not null;primary_key"`
	TodoID    int         `gorm:"type:int(20);not null;primary_key"`
	Body      string      `gorm:"type:varchar(200);not null"`
	LimitDate interface{} `gorm:"column:limit_date" sql:"null;type:timestamp"`
	CreatedAt time.Time   `gorm:"column:created_at" sql:"not null;type:timestamp"`
	UpdatedAt time.Time   `gorm:"column:updated_at" sql:"not null;type:timestamp"`
}

type Todos struct {
	Todos []Todo `json:"todos"`
}

func main() {
	// fmt.Print("\x1b[2J") // 画面全クリア
	fmt.Println("-----------------------------")
	fmt.Println(" Welcome to Go-Todo service.")
	fmt.Println(" Please choose menu.")

	for {
		fmt.Println("-----------------------------")
		fmt.Println(" [Login    ]: Enter")
		fmt.Println(" [User Menu]: 'u' & Enter")
		fmt.Println(" [Exit     ]: 'e' & Enter")
		choice := StdIn()

		if choice == "" {
			// directory = "/login"
			// Request()
			ShowLogin()
		} else if choice == "u" {
			// directory = "/user"
			// Request()
			ShowUserMenu()
		} else if choice == "e" {
			fmt.Println("-----------------------------")
			fmt.Println(" Goodbyeﾉｼ")
			break
		} else {
			fmt.Println(" Please enter again.")
		}
	}
}

func ShowLogin() {
	fmt.Println("-----------------------------")
	fmt.Println(" Please login to your account.")
	fmt.Println("-----------------------------")

	fmt.Printf(" [Name]:")
	name := StdIn()
	fmt.Printf(" [Pass]:")
	pass := StdIn()
	encryptedPass := ToSha1Hash(pass)

	// usercheckで確認
	directory = "/checkUser"
	value := url.Values{}
	value.Add("name", name)
	value.Add("password", encryptedPass)

	reqUrl := urlBase + directory

	if isExist, err := Request("POST", reqUrl, value); err != nil {
		fmt.Println("error isexist")
	} else {
		// userのidを返してもらう
		directory = "/returnUID"
		reqUrl := urlBase + directory
		uID, _ := Request("POST", reqUrl, value)
		uIDStr := StreamToString(uID.Body)
		uIDInt, _ := strconv.Atoi(uIDStr)
		Login(isExist.Body, uIDInt)
	}
}

func Login(isExistBody io.Reader, uID int) {
	isExist := StreamToString(isExistBody)
	// fmt.Println("isExist :" + isExist)
	if isExist == "ng" {
		fmt.Println("-----------------------------")
		fmt.Println(" Failed to login.")
		fmt.Println(" Back to main menu.")
	} else if isExist == "ok" {
		fmt.Println("-----------------------------")
		fmt.Println(" Login successed!")

		ShowTodo(uID)
	} else {

		// ここくっそ雑

		fmt.Println("-----------------------------")
		fmt.Println("errrrrrrrrrrrrrr")
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

func ShowUserMenu() {
	fmt.Println("-----------------------------")
	fmt.Println(" User manage menu")
	fmt.Println(" Please choose menu.")

	for {
		fmt.Println("-----------------------------")
		fmt.Println(" [Add user   ]: 'a' & Enter")
		fmt.Println(" [Delete user]: 'd' & Enter")
		fmt.Println(" [Exit       ]: 'e' & Enter")
		choice := StdIn()
		if choice == "a" {
			AddUser()
		} else if choice == "d" {
			DeleteUser()
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

func AddUser() {
	fmt.Printf(" Enter User name :")
	name := StdIn()
	fmt.Printf(" Enter Password  :")
	pass := StdIn()
	fmt.Printf(" Enter Password again:")
	passAgain := StdIn()
	if pass != passAgain {
		fmt.Println(" Passwords do not match.")
	} else {
		value := url.Values{}
		value.Add("name", name)
		value.Add("password", pass)
		directory = "/adduser"
		reqUrl := urlBase + directory

		if res, err := Request("POST", reqUrl, value); err != nil {
			fmt.Println("error add user")
		} else {
			resBody := StreamToString(res.Body)
			if resBody == "ng" {
				fmt.Println("")
				fmt.Println("-----------------------------")
				fmt.Println(" Add user failed.")
			} else if resBody == "ok" {
				fmt.Println("")
				fmt.Println("-----------------------------")
				fmt.Println(" Add user completed.")
			}
		}
	}
}

func DeleteUser() {
	fmt.Printf(" Enter User name want to delete :")
	name := StdIn()
	fmt.Println(" Are you sure?")
	fmt.Println(" [OK    ]: 'ok' & Enter")
	fmt.Println(" [CANCEL]: Enter")
	confirm := StdIn()
	if confirm != "ok" {
		fmt.Println(" User delete canceled.")
	} else {
		value := url.Values{}
		value.Add("name", name)
		directory = "/deleteuser"
		reqUrl := urlBase + directory

		if res, err := Request("POST", reqUrl, value); err != nil {
			fmt.Println("error delete user")
		} else {
			resBody := StreamToString(res.Body)
			if resBody == "ng" {
				fmt.Println("")
				fmt.Println("-----------------------------")
				fmt.Println(" Delete user failed.")
			} else if resBody == "ok" {
				fmt.Println("")
				fmt.Println("-----------------------------")
				fmt.Println(" Delete user completed.")
			}
		}
	}
}

func ShowTodo(uID int) {
	directory = "/checkTodo"
	reqUrl := urlBase + directory

	uIDStr := strconv.Itoa(uID)
	value := url.Values{}
	value.Add("id", uIDStr)

	if isExist, err := Request("POST", reqUrl, value); err != nil {
		fmt.Println("error isexist")
	} else {
		isExistBody := StreamToString(isExist.Body)
		// fmt.Println("isExist :" + isExistBody)
		if isExistBody == "ng" {
			fmt.Println("-----------------------------")
			fmt.Println(" Failed to get todos.")
			fmt.Println(" Back to main menu.")
		} else if isExistBody == "ok" {
			fmt.Println(" Here is your task.")
			// タスクの数を先に受け取っておいてその数だけループ
			// タスクのbodyとLimitdateを受け取る
			// タスクの何番目なのかを送ってその都度↑を受け取る

			directory = "/countTodo"
			reqUrl := urlBase + directory

			if length, err := Request("POST", reqUrl, value); err != nil {
				fmt.Println("error count todo")
			} else {
				lengthBody := StreamToString(length.Body)
				// fmt.Println(lengthBody)
				if lengthBody == "ng" {
					fmt.Println("error get task")
				} else if lengthBody == "zero" {
					fmt.Println("no task")
				} else {
					lengthNum, _ := strconv.Atoi(lengthBody)
					for i := 0; i < lengthNum; i++ {
						// var limitDate time.Time
						fmt.Println("-----------------------------")
						fmt.Println("ID: " + strconv.Itoa(i+1))

						value := url.Values{}
						value.Add("id", uIDStr)
						value.Add("num", strconv.Itoa(i))

						directory = "/getTodoBody"
						reqUrl := urlBase + directory

						if task, err := Request("POST", reqUrl, value); err != nil {
							fmt.Println("task get err")
						} else {
							fmt.Println("Task: " + StreamToString(task.Body))

							directory = "/getTodoLimitDate"
							reqUrl := urlBase + directory

							if limitDate, err := Request("POST", reqUrl, value); err != nil {
								fmt.Println("limitdate get err")
							} else {
								limitDateStr := StreamToString(limitDate.Body)
								// fmt.Println("limitDate:" + limitDateStr)
								if limitDateStr != "" {
									//layout := "Mon Jan _2 2006"
									layout := "2006-01-02 15:04:05"
									//limitDateParse := strings.TrimRight(limitDateStr, " +0000 UTC")
									limitDateParse := strings.Replace(limitDateStr, " +0000 UTC", "", 1)
									if t, err := time.Parse(layout, limitDateParse); err != nil {
										fmt.Printf("Failed : %s", err)
										fmt.Println()
									} else {
										if t.IsZero() != true {
											fmt.Println("Limit: " + t.Format("Mon Jan _2 2006"))
											// fmt.Println("Limit: " + t.Format("Mon Jan _2 2006"))
										}
									}
								}
							}
						}
					}

					fmt.Println("-----------------------------")
					fmt.Println("")
					fmt.Println(" [Add task    ]: Enter")
					fmt.Println(" [Delete task ]: 'd' & Enter")
					fmt.Println(" [Logout      ]: anykey & Enter")
					inputKey := StdIn()
					if inputKey == "d" {
						DeleteTask(uIDStr, lengthNum)
						// fmt.Println("delete")
					} else if inputKey == "" {
						AddTask(uIDStr)
						// fmt.Println("add")
					} else {
						//Logout()
						fmt.Println("logout")
					}
				}

			}
		}
	}
}

func AddTask(userID string) {
	var limitDateString string
	// layout := "2006-01-02 15:04:05"

	fmt.Printf("Task name: ")
	body := StdIn()

	fmt.Println(" [Set limitdate        ]: Enter")
	fmt.Println(" [Don't set limitdate  ]: anykey & Enter")
	fmt.Println("")
	addTask := StdIn()
	if addTask == "" {
		fmt.Printf(" Limit Year(ex. 2018年 → 2018): ")
		limitYear := StdIn()
		fmt.Printf(" Limit Month(ex. 1月 → 01): ")
		limitMonth := StdIn()
		fmt.Printf(" Limit Day(ex. 2日 → 02): ")
		limitDay := StdIn()
		limitDateString = limitYear + "-" + limitMonth + "-" + limitDay + " 00:00:00"
		// limitDate, _ = time.Parse(layout, limitDateString)
	} else {
		// limitDate = mysql.NullTime{Time: time.Now(), Valid: false}
		limitDateString = "0001-01-01 00:00:00"
	}

	value := url.Values{}
	value.Add("id", userID)
	value.Add("task", body)
	value.Add("limitDate", limitDateString)
	directory = "/addtask"
	reqUrl := urlBase + directory

	if res, err := Request("POST", reqUrl, value); err != nil {
		fmt.Println("error add todo")
	} else {
		resBody := StreamToString(res.Body)
		if resBody == "ng" {
			fmt.Println("")
			fmt.Println("-----------------------------")
			fmt.Println(" Add task failed.")
		} else if resBody == "ok" {
			fmt.Println("")
			fmt.Println("-----------------------------")
			fmt.Println(" Add task completed.")
		}
	}

}

func DeleteTask(userID string, length int) {
	fmt.Printf(" [Enter task ID to delete ]:")
	deleteIDString := StdIn()
	deleteID, _ := strconv.Atoi(deleteIDString)

	// userIDInt := userID.(int)
	// userIDInt, _ := strconv.Atoi(userID)
	// var todoID int
	if deleteID > length {
		fmt.Println(" Enter ID is wrong.")
	} else if deleteID <= 0 {
		fmt.Println(" Enter ID is wrong.")
	} else {
		// todoID = todos[deleteID-1].TodoID
		// model.DeleteTask(userIDInt, todoID)
		value := url.Values{}
		value.Add("id", userID)
		value.Add("num", deleteIDString)
		directory = "/deletetask"
		reqUrl := urlBase + directory

		if res, err := Request("POST", reqUrl, value); err != nil {
			fmt.Println("error delete todo")
		} else {
			resBody := StreamToString(res.Body)
			if resBody == "ng" {
				fmt.Println("")
				fmt.Println("-----------------------------")
				fmt.Println(" Delete task failed.")
			} else if resBody == "ok" {
				fmt.Println("")
				fmt.Println("-----------------------------")
				fmt.Println(" Delete task completed.")
			}
		}
	}
}

func Request(post string, reqUrl string, values url.Values) (*http.Response, error) {
	client := &http.Client{}
	body := strings.NewReader(values.Encode())

	if req, err := http.NewRequest(post, reqUrl, body); err != nil {
		return nil, err
	} else {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if resp, err := client.Do(req); err != nil {
			return nil, err
		} else {
			return resp, nil
		}
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

func StreamToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.String()
}
