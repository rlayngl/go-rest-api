package handlers

import (
	"encoding/json"
	"net/http"
	"restapi/internal/database"
	"restapi/internal/models"
	"strconv"
	"strings"
)

type Handler struct {
	store *database.TaskStore
}

func NewHandler(store *database.TaskStore) *Handler {
	return &Handler{store: store}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.GetAll()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, tasks)
}

func (h *Handler) GetID(w http.ResponseWriter, r *http.Request) (int, error) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idString := pathParts[0]

	id, idErr := strconv.Atoi(idString)
	if idErr != nil {
		respondWithError(w, http.StatusBadRequest, "invalid task id")
		return 0, idErr
	}
	return id, nil
}

func (h *Handler) GetTaskById(w http.ResponseWriter, r *http.Request) {
	id, idErr := h.GetID(w, r)
	if idErr != nil {
		respondWithError(w, http.StatusBadRequest, idErr.Error())
		return
	}

	task, err := h.store.GetByID(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, task)
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask models.CreateTaskInput

	if decodingErr := json.NewDecoder(r.Body).Decode(&newTask); decodingErr != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(newTask.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "no request body")
		return
	}
	task, err := h.store.Create(&newTask)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, task)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, idErr := h.GetID(w, r)

	if idErr != nil {
		respondWithError(w, http.StatusBadRequest, idErr.Error())
		return
	}

	var newTask models.UpdateTaskInput
	if decodingErr := json.NewDecoder(r.Body).Decode(&newTask); decodingErr != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if newTask.Title != nil && strings.TrimSpace(*newTask.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "no request body")
		return
	}

	task, err := h.store.Update(id, &newTask)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, task)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, idErr := h.GetID(w, r)
	if idErr != nil {
		respondWithError(w, http.StatusBadRequest, idErr.Error())
		return
	}

	if err := h.store.Delete(id); err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "success"})
}
