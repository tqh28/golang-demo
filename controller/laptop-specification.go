package controller

import (
	"log"
	"net/http"

	"github.com/kataras/iris"

	"github.com/go-xorm/xorm"

	"encoding/json"
)

type LaptopSpecification struct {
	Laptop        `xorm:"extends"`
	Specification `xorm:"extends"`
}

func DeclareLaptopSpecification(application *iris.Application, engine *xorm.Engine) {

	application.Get("/laptop-specification/get-all", func(ctx iris.Context) {
		var message string
		var statusCode int
		var laptopSpecification []LaptopSpecification

		err := engine.Table("laptop").Join("INNER", "specification", "laptop.LaptopID = specification.LaptopID").Find(&laptopSpecification)

		if err != nil {
			log.Printf("Laptop specification join query error: %v", err)
			statusCode = http.StatusInternalServerError
			message = err.Error()
		} else {
			res, err := json.Marshal(laptopSpecification)
			if err != nil {
				log.Printf("Laptop specification parse data error: %v", err)
				statusCode = http.StatusInternalServerError
				message = err.Error()
			} else {
				message = string(res)
				statusCode = http.StatusOK
			}
		}

		ctx.StatusCode(statusCode)
		ctx.Writef(message)
	})

}
