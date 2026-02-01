package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"tasks-api/internal/models"
	"tasks-api/internal/storage"
	"tasks-api/internal/storage/memory"
)

type Handler struct{ Store storage.Storage }

func New(s storage.Storage) *Handler { return &Handler{Store: s} }

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

// /tasks (GET, POST)
func (h *Handler) TasksCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tasks := h.Store.List()
		writeJSON(w, http.StatusOK, tasks)
		return

	case http.MethodPost:
		var in struct {
			Title     string `json:"title"`
			Done      bool   `json:"done"`
			CreatedAt string `json:"created_at,omitempty"`
		}

		if err := decodeJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		in.Title = strings.TrimSpace(in.Title)
		if in.Title == "" {
			writeError(w, http.StatusBadRequest, "title is required")
			return
		}

		created, err := h.Store.Create(models.Task{
			Title:     in.Title,
			Done:      in.Done,
			CreatedAt: in.CreatedAt,
		})
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		writeJSON(w, http.StatusCreated, created)
		return

	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
}

// /tasks/{id} (GET, PUT, DELETE)
func (h *Handler) TaskItem(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r.URL.Path)
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		task, found := h.Store.Get(id)
		if !found {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeJSON(w, http.StatusOK, task)
		return

	case http.MethodPut:
		var in struct {
			Title     string `json:"title"`
			Done      bool   `json:"done"`
			CreatedAt string `json:"created_at,omitempty"` // игнорируем при обновлении
		}

		if err := decodeJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		in.Title = strings.TrimSpace(in.Title)
		if in.Title == "" {
			writeError(w, http.StatusBadRequest, "title is required")
			return
		}

		updated, err := h.Store.Update(id, models.Task{
			Title: in.Title,
			Done:  in.Done,
		})
		if err != nil {
			if errors.Is(err, memory.ErrNotFound) {
				writeError(w, http.StatusNotFound, "task not found")
				return
			}
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		writeJSON(w, http.StatusOK, updated)
		return

	case http.MethodDelete:
		if err := h.Store.Delete(id); err != nil {
			if errors.Is(err, memory.ErrNotFound) {
				writeError(w, http.StatusNotFound, "task not found")
				return
			}
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		// 204 — без тела
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNoContent)
		return

	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
}

// ожидаем строго /tasks/{id} без лишних сегментов
func parseID(path string) (int, bool) {
	if !strings.HasPrefix(path, "/tasks/") {
		return 0, false
	}

	rest := strings.TrimPrefix(path, "/tasks/")
	rest = strings.Trim(rest, "/")

	if rest == "" {
		return 0, false
	}
	if strings.Contains(rest, "/") {
		return 0, false
	}

	id, err := strconv.Atoi(rest)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}
