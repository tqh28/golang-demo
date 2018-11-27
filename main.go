package main

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris"

	"demo/controller"

	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
)

const (
	privKeyPath = "keys/jwtRS256.key"
	pubKeyPath  = "keys/jwtRS256.key.pub"
)

func initKeys() (*rsa.PrivateKey, *rsa.PublicKey) {
	signKey, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatal("Reading private key error: %s", err)
	}
	signKeyObject, err := jwt.ParseRSAPrivateKeyFromPEM(signKey)
	if err != nil {
		log.Fatal("Convert sign key error: %s", err)
	}

	verifyKey, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatal("Error reading public key: %s", err)
	}
	verifyKeyObject, err := jwt.ParseRSAPublicKeyFromPEM(verifyKey)
	if err != nil {
		log.Fatal("Convert verify key error: %s", err)
	}

	return signKeyObject, verifyKeyObject
}

func main() {
	application := iris.New()

	// open data base connection
	engine, err := xorm.NewEngine("mssql", "server=localhost;user id=sa;password=qwerty@123;database=demo")
	if err != nil {
		// log and exit programe
		application.Logger().Fatal("Connect to database error: %v", err)
	}

	iris.RegisterOnInterrupt(func() {
		engine.Close()
	})
	engine.SetTableMapper(core.SameMapper{}) // set mapper for table column
	engine.ShowSQL(true)                     // debug mode

	signKey, verifyKey := initKeys()

	jwtHandler := jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})
	application.Use(func(ctx iris.Context) {
		// ignore login require jwt
		if strings.Contains(ctx.Path(), "login") {
			ctx.Next()
		} else {
			jwtHandler.Serve(ctx)
		}
	})

	controller.DeclareLogin(application, engine, signKey)
	controller.DeclareLaptop(application, engine)
	controller.DeclareSpecification(application, engine)
	controller.DeclareLaptopSpecification(application, engine)
	controller.DeclareFoward(application)
	controller.DeclareUpload(application)
	controller.DeclareError(application)
	controller.DeclareAdmin(application)

	application.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
