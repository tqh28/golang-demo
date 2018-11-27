package controller

import (
	"log"

	"github.com/kataras/iris"
)

func DeclareError(application *iris.Application) {

	application.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		log.Printf("User request url: %s", ctx.Path())
		ctx.Writef("API not found")
	})

}
