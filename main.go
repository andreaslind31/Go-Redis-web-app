package main

import (
	"net/http"

	"github.com/andreaslind31/Go-Redis-web-app/models"
	"github.com/andreaslind31/Go-Redis-web-app/routes"
	"github.com/andreaslind31/Go-Redis-web-app/utils"
)

func main() {
	models.InitializeDb()
	utils.LoadTemplates("templates/*.html")
	r := routes.NewRouter()
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
