package main

import (
	"fmt"
	"net/http"
	"yboost-portfolio/src/helper"
	"yboost-portfolio/src/routes"
	"os"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	// Chargement des variables d'environnement
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erreur lors du chargement des variables d'environnement: %v", err)
	}
	// Chargement des templates
	helper.Load()
	// Chargement des routes du serveur
	serveRouter := routes.MainRouter()

	port := os.Getenv("PORT")
	// Message d'information indiquant que le serveur est lancé
	fmt.Println("Server started")
	// Lancement du serveur HTTP sur le port 8081
	http.ListenAndServe(":"+ port, serveRouter)
}
