package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

func main() {
	server := gin.Default()

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
		fileExt := filepath.Ext(file.Filename)

		// save the file
		fileId := ksuid.New()
		fileName := fmt.Sprintf("gpfile%s", fileId.String())
		gpFileName := fileName + fileExt
		c.SaveUploadedFile(file, gpFileName)

		// create midi file
		err = exec.Command("./GuitarProToMidi", gpFileName).Run()
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "file could not be converted",
			})
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
