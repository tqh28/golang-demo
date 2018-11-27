package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kataras/iris"
)

func DeclareFoward(application *iris.Application) {

	application.Get("/foward", func(ctx iris.Context) {
		var message string
		var statusCode int

		url := fmt.Sprint("http://localhost:8080/")

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println("Create request error: ", err)
			message = err.Error()
			statusCode = http.StatusInternalServerError
		} else {
			client := &http.Client{}
			res, err := client.Do(req)
			if err != nil {
				log.Println("Run request error: ", err)
				message = err.Error()
				statusCode = http.StatusInternalServerError
			}

			defer res.Body.Close()
			bodyBytes, err := ioutil.ReadAll(res.Body)
			if err != nil {
				message = err.Error()
				statusCode = http.StatusInternalServerError
			} else {
				message = string(bodyBytes)
				statusCode = res.StatusCode
			}
		}

		ctx.StatusCode(statusCode)
		ctx.Writef(message)
	})

}
