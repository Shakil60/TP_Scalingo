package routes

import (
	"net/http"
	"yboost-portfolio/src/controllers"
)

// MainRouter initialise et retourne le routeur principal de l'application
func MainRouter() *http.ServeMux {

	// Création du routeur principal
	mainRouter := http.NewServeMux()

	// Route pour la page d'accueil
	mainRouter.HandleFunc("/", controllers.HomeDisplay)

	// Route pour la page de todo
	mainRouter.HandleFunc("/todo_list", controllers.TodoListDisplay)

	// Configuration du serveur de fichiers statiques (CSS, images, etc.)
	// On pointe sur le dossier "assets" à la racine du projet
	fileServerHandler := http.FileServer(http.Dir("assets"))

	// Route permettant de servir les fichiers statiques via /static/
	mainRouter.Handle("/static/", http.StripPrefix("/static/", fileServerHandler))

	return mainRouter
}
