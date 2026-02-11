package models

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Todo représente une tâche à faire
// Elle est stockée dans la table SQL `todo` (id, title, completed)
type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"` // mappé sur la colonne `title`
	Done bool   `json:"done"` // mappé sur la colonne `completed`
}

// TodoStore gère le stockage des todos en base de données
type TodoStore struct {
	db *sql.DB
}

// NewTodoStore crée une nouvelle instance de TodoStore
// et ouvre une connexion à la base de données
func NewTodoStore() *TodoStore {
	db := getDB()
	return &TodoStore{
		db: db,
	}
}

// GetAll retourne toutes les todos depuis la base de données
func (ts *TodoStore) GetAll() []Todo {
	rows, err := ts.db.Query("SELECT id, title, completed FROM todo")
	if err != nil {
		log.Printf("Erreur lors de la récupération des todos: %v", err)
		return []Todo{}
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Text, &t.Done); err != nil {
			log.Printf("Erreur lors du scan d'une todo: %v", err)
			continue
		}
		todos = append(todos, t)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Erreur lors de l'itération des todos: %v", err)
	}

	return todos
}

// Add ajoute une nouvelle todo dans la table `todo`
func (ts *TodoStore) Add(text string) Todo {
	result, err := ts.db.Exec(
		"INSERT INTO todo (title, completed) VALUES (?, ?)",
		text,
		false,
	)
	if err != nil {
		log.Printf("Erreur lors de l'ajout de la todo: %v", err)
		return Todo{}
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Erreur lors de la récupération de l'ID de la todo: %v", err)
		return Todo{}
	}

	return Todo{
		ID:   int(lastID),
		Text: text,
		Done: false,
	}
}

// Delete supprime une todo par son ID dans la base de données
func (ts *TodoStore) Delete(id int) bool {
	result, err := ts.db.Exec("DELETE FROM todo WHERE id = ?", id)
	if err != nil {
		log.Printf("Erreur lors de la suppression de la todo: %v", err)
		return false
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Erreur lors de la récupération du nombre de lignes supprimées: %v", err)
		return false
	}

	return rowsAffected > 0
}

// Toggle change l'état (fait/non fait) d'une todo en base de données
func (ts *TodoStore) Toggle(id int) bool {
	result, err := ts.db.Exec(
		"UPDATE todo SET completed = NOT completed WHERE id = ?",
		id,
	)
	if err != nil {
		log.Printf("Erreur lors du changement d'état de la todo: %v", err)
		return false
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Erreur lors de la récupération du nombre de lignes modifiées: %v", err)
		return false
	}

	return rowsAffected > 0
}

func getDB() *sql.DB {
	user := os.Getenv("SCALINGO_MYSQL_USER")
	password := os.Getenv("SCALINGO_MYSQL_PASSWORD")
	host := os.Getenv("SCALINGO_MYSQL_HOST")
	port := os.Getenv("SCALINGO_MYSQL_PORT")
	dbname := os.Getenv("SCALINGO_MYSQL_DBNAME")
	
	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Erreur lors de la connexion à la base de données: %v", err)
	}
	return db
}

// Instance globale du store
var GlobalTodoStore = NewTodoStore()
