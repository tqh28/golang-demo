package controller

import (
	"log"
	"net/http"

	"github.com/kataras/iris"

	"github.com/go-xorm/xorm"
)

type Specification struct {
	SpecificationID int32  `xorm:"pk not null autoincr 'specificationID'"`
	LaptopID        int32  `xorm:"not null 'laptopID'"`
	Dimension       string `xorm:"not null 'dimension'"`
	Weight          string `xorm:"'weight'"`
}

func DeclareSpecification(application *iris.Application, engine *xorm.Engine) {

	err := engine.Sync2(new(Specification))

	if err != nil {
		log.Fatalf("orm failed to initialized product table: %v", err)
	}

	application.Post("/specification/insert", func(ctx iris.Context) {
		var message string
		var statusCode int

		laptopID := int32(ctx.PostValueInt64Default("laptop-id", -1))

		spec := &Specification{LaptopID: laptopID, Dimension: ctx.PostValue("dimension"), Weight: ctx.PostValue("weight")}
		res, err := engine.Insert(spec)
		if err != nil {
			log.Printf("Insert specification error: %v", err)
			statusCode = http.StatusInternalServerError
			message = err.Error()
		} else if res == 1 {
			statusCode = http.StatusOK
			message = ""
		} else {
			statusCode = http.StatusInternalServerError
		}

		ctx.StatusCode(statusCode)
		ctx.Writef(message)
	})

}
