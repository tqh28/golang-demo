package controller

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/kataras/iris"
)

const uploadsDir = "./public/uploads/"

func DeclareUpload(application *iris.Application) {
	application.Post("/upload", func(ctx iris.Context) {
		var message string
		var statusCode int

		file, _, err := ctx.FormFile("file")

		if err != nil {
			log.Printf("Upload file error: %v", err)
			statusCode = http.StatusBadRequest
			message = err.Error()
		} else {
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}

			if err := scanner.Err(); err != nil {
				log.Printf("Scan file error: %v", err)
				message = err.Error()
				statusCode = http.StatusInternalServerError
			} else {
				message = "Successful"
				statusCode = http.StatusOK
			}
		}
		ctx.StatusCode(statusCode)
		ctx.Writef(message)
	})

	application.Post("/multi-upload", func(ctx iris.Context) {
		// Get the max post value size passed via iris.WithPostMaxMemory.
		maxSize := ctx.Application().ConfigurationReadOnly().GetPostMaxMemory()

		err := ctx.Request().ParseMultipartForm(maxSize)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString(err.Error())
			return
		}

		form := ctx.Request().MultipartForm

		files := form.File["files[]"]
		failures := 0
		for _, file := range files {

			src, err := file.Open()
			if err != nil {
				return
			}
			defer src.Close()

			buf := make([]byte, 4)
			if _, err := io.ReadFull(src, buf); err != nil {
				log.Println(err)
			}
			fmt.Printf("%s\n", buf)

			// minimal read size bigger than io.Reader stream
			longBuf := make([]byte, 64)
			if _, err := io.ReadFull(src, longBuf); err != nil {
				fmt.Println("error:", err)
			}
			fmt.Printf("%s\n", longBuf)

			// _, err = saveUploadedFile(file, "./uploads")
			// if err != nil {
			// 	failures++
			// 	ctx.Writef("failed to upload: %s\n", file.Filename)
			// }
		}
		ctx.Writef("%d files uploaded", len(files)-failures)

	})
}
