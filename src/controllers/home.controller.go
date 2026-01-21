package controllers

import (
	"net/http"
	"strconv"
	"yboost-portfolio/src/helper"
	"yboost-portfolio/src/models"
)

// HomeDisplay est le contrôleur pour la page d'accueil
func HomeDisplay(w http.ResponseWriter, r *http.Request) {
	// Pas de données spécifiques nécessaires pour la page d'accueil
	helper.RenderTemplate(w, r, "index", nil)
}

// TodoListDisplay est le contrôleur pour la page de todo list
func TodoListDisplay(w http.ResponseWriter, r *http.Request) {
	// Gestion des méthodes POST pour les actions
	if r.Method == http.MethodPost {
		action := r.FormValue("action")

		switch action {
		case "add":
			// Ajouter une nouvelle todo
			text := r.FormValue("text")
			if text != "" {
				models.GlobalTodoStore.Add(text)
			}
			http.Redirect(w, r, "/todo_list", http.StatusSeeOther)
			return

		case "delete":
			// Supprimer une todo
			idStr := r.FormValue("id")
			if id, err := strconv.Atoi(idStr); err == nil {
				models.GlobalTodoStore.Delete(id)
			}
			http.Redirect(w, r, "/todo_list", http.StatusSeeOther)
			return

		case "toggle":
			// Cocher/décocher une todo
			idStr := r.FormValue("id")
			if id, err := strconv.Atoi(idStr); err == nil {
				models.GlobalTodoStore.Toggle(id)
			}
			http.Redirect(w, r, "/todo_list", http.StatusSeeOther)
			return
		}
	}

	// Récupération de toutes les todos pour l'affichage
	todos := models.GlobalTodoStore.GetAll()

	// Données à passer au template
	data := map[string]interface{}{
		"Todos": todos,
	}

	helper.RenderTemplate(w, r, "todo_list", data)
}
