package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

// User JSON format
type User struct {
	ID          int64  `json:"id" binding:"required"`
	FullName    string `json:"fullName" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	ActiveEmail string `json:"activeEmail" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

func main() {
	router := gin.Default()
	router.GET("/users", getUsers)
	router.POST("/user", saveUser)
	router.PUT("/user", updateUser)
	router.DELETE("/user", deleteUser)
	router.GET("/user/:id", getUserByIDPath)
	router.Run(":8080")
}

func connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:hello1234@tcp(127.0.0.1:3306)/abhi")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getUsers(c *gin.Context) {
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()
	rows, errDBConn := db.Query("SELECT * FROM master_user")
	if errDBConn != nil {
		fmt.Println(errDBConn.Error())
		return
	}
	var users []User
	var newUser = User{}
	for rows.Next() {

		var err = rows.Scan(&newUser.ID, &newUser.FullName, &newUser.PhoneNumber, &newUser.ActiveEmail, &newUser.Password)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		users = append(users, newUser)
	}
	rows.Close()
	c.JSON(200, users)
}

func saveUser(c *gin.Context) {
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var user User
	c.ShouldBindJSON(&user)
	res, errSave := db.Exec("INSERT INTO master_user VALUES (?, ?, ?, ?, ?)", user.ID, user.FullName, user.PhoneNumber, user.ActiveEmail, user.Password)
	if errSave != nil {
		fmt.Println(errSave.Error())
		return
	}
	defer db.Close()
	lastID, _ := res.LastInsertId()

	user.ID = lastID
	c.JSON(200, user)
}

func updateUser(c *gin.Context) {
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var user User
	c.ShouldBindJSON(&user)
	_, errSave := db.Exec("UPDATE master_user SET full_name = ? WHERE id = ?", user.FullName, user.ID)
	if errSave != nil {
		fmt.Println(errSave.Error())
		return
	}
	defer db.Close()
	c.JSON(200, user)
}

func deleteUser(c *gin.Context) {
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var user User
	c.ShouldBindJSON(&user)
	_, errSave := db.Exec("DELETE FROM master_user WHERE id = ?", user.ID)
	if errSave != nil {
		fmt.Println(errSave.Error())
		return
	}
	defer db.Close()
	c.JSON(200, user)
}

func getUserByIDPath(c *gin.Context) {
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()
	id := c.Param("id")
	rows, err := db.Query("SELECT * FROM master_user WHERE id = ?", id)
	var users []User
	var newUser = User{}
	for rows.Next() {

		var err = rows.Scan(&newUser.ID, &newUser.FullName, &newUser.PhoneNumber, &newUser.ActiveEmail, &newUser.Password)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		users = append(users, newUser)
	}
	rows.Close()
	c.JSON(200, users)
}
