package models

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

// Todo représente une tâche à faire
type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

// TodoStore gère le stockage des todos en mémoire
type TodoStore struct {
	todos    []Todo
	mu       sync.RWMutex
	nextID   int
	filename string
}

// NewTodoStore crée une nouvelle instance de TodoStore
func NewTodoStore() *TodoStore {
	store := &TodoStore{
		todos:    make([]Todo, 0),
		nextID:   1,
		filename: "../todo.json",
	}
	store.load()
	return store
}

// load charge les todos depuis le fichier JSON
func (ts *TodoStore) load() {
	file, err := os.ReadFile(ts.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return // Le fichier n'existe pas encore, on commence vide
		}
		log.Printf("Erreur lors de la lecture du fichier todo: %v", err)
		return
	}

	err = json.Unmarshal(file, &ts.todos)
	if err != nil {
		log.Printf("Erreur lors du décodage du fichier todo: %v", err)
		return
	}

	// Mettre à jour nextID basé sur l'ID le plus élevé
	maxID := 0
	for _, todo := range ts.todos {
		if todo.ID > maxID {
			maxID = todo.ID
		}
	}
	ts.nextID = maxID + 1
}

// save sauvegarde les todos dans le fichier JSON
func (ts *TodoStore) save() {
	data, err := json.MarshalIndent(ts.todos, "", "  ")
	if err != nil {
		log.Printf("Erreur lors de l'encodage des todos: %v", err)
		return
	}

	err = os.WriteFile(ts.filename, data, 0644)
	if err != nil {
		log.Printf("Erreur lors de l'écriture du fichier todo: %v", err)
	}
}

// GetAll retourne toutes les todos
func (ts *TodoStore) GetAll() []Todo {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	// Créer une copie pour éviter les problèmes de concurrence
	result := make([]Todo, len(ts.todos))
	copy(result, ts.todos)
	return result
}

// Add ajoute une nouvelle todo
func (ts *TodoStore) Add(text string) Todo {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	todo := Todo{
		ID:   ts.nextID,
		Text: text,
		Done: false,
	}
	ts.todos = append(ts.todos, todo)
	ts.nextID++
	ts.save()
	return todo
}

// Delete supprime une todo par son ID
func (ts *TodoStore) Delete(id int) bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	for i, todo := range ts.todos {
		if todo.ID == id {
			ts.todos = append(ts.todos[:i], ts.todos[i+1:]...)
			ts.save()
			return true
		}
	}
	return false
}

// Toggle change l'état (fait/non fait) d'une todo
func (ts *TodoStore) Toggle(id int) bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	for i := range ts.todos {
		if ts.todos[i].ID == id {
			ts.todos[i].Done = !ts.todos[i].Done
			ts.save()
			return true
		}
	}
	return false
}

// Instance globale du store
var GlobalTodoStore = NewTodoStore()
