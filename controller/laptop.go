package controller

import (
	"encoding/json"
	"net/http"

	"github.com/kataras/iris"

	"log"

	"github.com/go-xorm/xorm"
)

type Laptop struct {
	LaptopID    int32  `xorm:"pk not null autoincr 'laptopID'"`
	Name        string `xorm:"not null 'name'"`
	Description string `xorm:"'description'"`
}

func DeclareLaptop(application *iris.Application, engine *xorm.Engine) {

	err := engine.Sync2(new(Laptop))

	if err != nil {
		log.Fatalf("orm failed to initialized laptop table: %v", err)
	}

	application.Post("/laptop/insert/", func(ctx iris.Context) {
		var message string
		var statusCode int
		lap := &Laptop{Name: ctx.PostValue("name"), Description: ctx.PostValue("desc")}
		res, err := engine.Insert(lap)
		if err == nil && res > 0 {
			statusCode = http.StatusOK
			message = "Success"
		} else {
			statusCode = http.StatusInternalServerError
			message = err.Error()
		}

		ctx.StatusCode(statusCode)
		ctx.Writef(message)
	})

	application.Get("/laptop/get-all", func(ctx iris.Context) {
		var message string
		var statusCode int

		var allLaptop []Laptop
		err := engine.Find(&allLaptop)
		if err != nil {
			log.Printf("Query all laptop error: %v", err)
			statusCode = http.StatusInternalServerError
			message = err.Error()
		} else {
			res, err := json.Marshal(allLaptop)
			if err != nil {
				statusCode = http.StatusInternalServerError
				message = err.Error()
			} else {
				statusCode = http.StatusOK
				message = string(res)
			}
		}

		ctx.StatusCode(statusCode)
		ctx.Writef(message)
	})

}
