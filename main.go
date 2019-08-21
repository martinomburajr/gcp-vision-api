package main

import (
	"fmt"
	"github.com/martinomburajr/gcp-vision-api/app"
	"github.com/martinomburajr/gcp-vision-api/app/router"
	"log"
	"net/http"
	"os"
	"strconv"
)

var PORT = 8080

func init() {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", app.CredentialsLocalPath)

	if err != nil {
		msg := "error setting GOOGLE_APPLICATION_CREDENTIALS."
		log.Fatalf(msg)
	}
}

func main() {
	server := router.GCPVisionAPIServer{}
	err := server.Init()
	if err != nil {
		log.Fatalf("error initializing server | %s ", err.Error())
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Printf("error parsing ENV: PORT defaulting to %d ", PORT)
	} else {
		envPort, err := strconv.ParseInt(port, 10, 64)
		if err != nil {
			log.Printf("error parsing ENV: PORT defaulting to %d | %v", PORT, err)
		}
		PORT = int(envPort)
	}

	log.Printf(fmt.Sprintf(" GCP-VISION-API Listening on PORT %d", PORT))

	err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), server.Routes())
	if err != nil {
		log.Println(err)
	}
}
