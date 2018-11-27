package controller

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID   int32  `xorm:"pk not null autoincr 'userID'"`
	Username string `xorm:"not null 'username'"`
	Password string `xorm:"not null 'password'"`
	Blocked  bool   `xorm:"not null 'blocked'"`
	Disabled bool   `xorm:"not null 'disabled'"`
	Deleted  bool   `xorm:"not null 'deleted'"`
}

func DeclareLogin(application *iris.Application, engine *xorm.Engine, signKey *rsa.PrivateKey) {

	err := engine.Sync2(new(User))
	if err != nil {
		log.Fatalf("orm failed to initialized user table: %v", err)
	}

	application.Post("/login", func(ctx iris.Context) {
		var message string
		var statusCode int
		username := ctx.PostValue("username")
		password := ctx.PostValue("password")

		res, loginFailMess := checkLogin(engine, username, password)

		if res {
			signer := jwt.New(jwt.SigningMethodRS256)
			claims := make(jwt.MapClaims)

			claims["username"] = username
			claims["role"] = "all"
			signer.Claims = claims

			tokenString, err := signer.SignedString(signKey)
			if err != nil {
				log.Printf("Sign token error: %v", err)
				message = err.Error()
				statusCode = http.StatusInternalServerError
			} else {
				message = fmt.Sprintf("Bearer %s", tokenString)
				statusCode = http.StatusOK
			}
		} else {
			message = loginFailMess
			statusCode = http.StatusBadRequest
		}

		ctx.StatusCode(statusCode)
		ctx.Writef(message)
	})
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkLogin(engine *xorm.Engine, username, password string) (bool, string) {
	var result bool
	var message string
	user := User{Username: username}

	has, err := engine.Get(&user)
	if err != nil {
		result, message = false, err.Error()
	} else if !has || user.Deleted {
		result, message = false, "User is not existed"
	} else {
		passwordCompareError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

		if passwordCompareError == nil {
			if user.Blocked {
				result, message = false, "User is blocked"
			} else {
				result, message = true, ""
			}
		} else {
			result, message = false, "Password is not correct"
		}
	}

	return result, message
}

func NewUser(engine *xorm.Engine, username, password string) (bool, string) {
	var result bool
	var message string
	hashedPassword, err := hashPassword(password)
	if err != nil {
		result, message = false, err.Error()
	} else {
		user := User{
			Username: username,
			Password: hashedPassword,
			Blocked:  false,
			Disabled: false,
			Deleted:  false,
		}
		res, err := engine.Insert(&user)

		if err == nil && res == 1 {
			result, message = true, ""
		} else {
			result, message = false, err.Error()
		}
	}
	return result, message
}

func BlockUser(engine *xorm.Engine, username string) {

}

func DeleteUser(engine *xorm.Engine, username string) {

}
