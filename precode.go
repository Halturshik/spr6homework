package main

import (
	"encoding/json"
	"fmt"
	"io"
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
		Description: "Протестировать финальное задание с помощью Postman",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, "Ошибка при обработке данных", http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func AddNewTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка при обработке данных", http.StatusInternalServerError)
		return
	}

	var newTask Task
	if err := json.Unmarshal(body, &newTask); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if newTask.Description == "" || len(newTask.Applications) == 0 {
		http.Error(w, "Отсутствуют обязательные поля", http.StatusBadRequest)
		return
	}

	newTask.ID = fmt.Sprintf("%d", len(tasks)+1)

	tasks[newTask.ID] = newTask

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response, _ := json.Marshal(newTask)
	w.Write(response)
}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, found := tasks[id]
	if !found {
		http.Error(w, "Такой задачи не существует", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "Ошибка при обработке данных", http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func DeleteTaskByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, found := tasks[id]; !found {
		http.Error(w, "Такой задачи не существует", http.StatusBadRequest)
		return
	}

	delete(tasks, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := []byte(`"Задача успешно удалена"`)
	w.Write(response)
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", GetAllTasks)
	r.Post("/tasks", AddNewTask)
	r.Get("/tasks/{id}", GetTaskByID)
	r.Delete("/tasks/{id}", DeleteTaskByID)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
