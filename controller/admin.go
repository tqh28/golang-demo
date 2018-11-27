package controller

import (
	"log"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
)

func DeclareAdmin(application *iris.Application) {

	adminRoutes := application.Party("/admin", func(ctx iris.Context) {
		log.Println("pre admin log")
		ctx.Next()
	})

	adminRoutes.Get("/profile", func(ctx iris.Context) {
		user := ctx.Values().Get("jwt").(*jwt.Token)

		ctx.Writef("This is an authenticated request\n")
		ctx.Writef("Claim content:\n")
		ctx.Writef("%s", user.Claims)
	})

}
