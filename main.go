package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

const (
	maxFileUploadSize int64  = 500_000 // 500kb
	fileUploadType    string = "application/octet-stream"
)

func main() {
	server := gin.Default()

	// CORS config for local development
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:3000",
		"https://fast-river-62884.herokuapp.com/",
	}
	server.Use(cors.New(corsConfig))

	// Routes
	server.POST("/compute", func(c *gin.Context) {
		// get file
		file, err := c.FormFile("gpfile")
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusUnsupportedMediaType, gin.H{
				"error": "payload must be type 'multipart/form-data'",
			})
			return
		}

		// check file type
		contentType := file.Header["Content-Type"][0]
		fmt.Println(contentType)
		if contentType != fileUploadType {
			log.Println("unsupported media type:", contentType)
			c.JSON(http.StatusUnsupportedMediaType, gin.H{
				"error": "file type must be " + fileUploadType,
			})
			return
		}

		// check file size
		if file.Size > maxFileUploadSize {
			log.Println("file upload size too large:", file.Size)
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "file upload size too large",
			})
			return
		}

		// load gp2midi binary
		prgPath, err := filepath.Abs("./GuitarProToMidi")
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "unable to get absolute filepath of gp2midi",
			})
			return
		}

		// save the file
		fileId := ksuid.New()
		fileName := fmt.Sprintf("gpfile%s", fileId.String())
		fileExt := filepath.Ext(file.Filename)
		gpFileName := fileName + fileExt
		c.SaveUploadedFile(file, gpFileName)

		// create midi file
		err = exec.Command(prgPath, gpFileName).Run()
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "file could not be converted",
			})
			os.Remove(gpFileName)
			return
		}

		// return midi file
		midiFileName := fileName + ".mid"
		c.File(midiFileName)

		// delete old files
		os.Remove(gpFileName)
		os.Remove(midiFileName)

		log.Printf("success: %s --> %s\n", gpFileName, midiFileName)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8229"
	}

	server.Run(":" + port)
}
