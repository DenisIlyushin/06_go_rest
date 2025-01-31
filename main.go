package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// sendJSONResponse Формирует ответ, с указанным статусом и JSON в теле ответа
func sendJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Ошибка кодирования JSON", http.StatusInternalServerError)
	}
}

// createTask Создает новое задание
func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Ошибка декодирования JSON", http.StatusBadRequest)
		return
	}
	if _, exists := tasks[task.ID]; exists {
		http.Error(w, fmt.Sprintf("Задача с id %v уже существует", task.ID), http.StatusConflict)
		return
	}
	tasks[task.ID] = task
	sendJSONResponse(w, http.StatusCreated, task)
}

// getAllTasks Возвращает список заданий
func getAllTasks(w http.ResponseWriter, r *http.Request) {
	sendJSONResponse(w, http.StatusOK, tasks)
}

// getTaskByID Возвращает задание с указанным ID
func getTaskById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, exists := tasks[id]
	if !exists {
		http.Error(w, fmt.Sprintf("Задачи с id %v не существует", id), http.StatusNotFound)
		return
	}
	sendJSONResponse(w, http.StatusOK, task)
}

// getTaskByID Удаляет задание с указанным ID
func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, exists := tasks[id]
	if exists {
		delete(tasks, id)
	}
	if !exists {
		http.Error(w, fmt.Sprintf("Задачи с id %v не существует", id), http.StatusNotFound)
		return
	}
	sendJSONResponse(w, http.StatusOK, task)
}

func main() {
	r := chi.NewRouter()

	r.Post("/tasks", createTask)
	r.Get("/tasks", getAllTasks)
	r.Get("/tasks/{id}", getTaskById)
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
