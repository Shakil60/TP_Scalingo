package main

import (
	"fmt"
	"net/http"
	"yboost-portfolio/helper"
	"yboost-portfolio/routes"
)

func main() {
	// Chargement des templates
	helper.Load()
	// Chargement des routes du serveur
	serveRouter := routes.MainRouter()
	// Message d'information indiquant que le serveur est lancé
	fmt.Println("Server started at http://localhost:8080")
	// Lancement du serveur HTTP sur le port 8081
	http.ListenAndServe("localhost:8080", serveRouter)
}
