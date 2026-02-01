package main

import (
	"log"
	"net/http"

	"tasks-api/internal/handlers"
	httpmw "tasks-api/internal/http"
	"tasks-api/internal/storage/memory"
)

func main() {
	store := memory.New()
	h := handlers.New(store)

	mux := http.NewServeMux()

	// PRO: /health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte("{\"error\":\"method not allowed\"}\n"))
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{\"status\":\"ok\"}\n"))
	})

	mux.HandleFunc("/tasks", h.TasksCollection) // GET, POST
	mux.HandleFunc("/tasks/", h.TaskItem)       // GET, PUT, DELETE

	handler := httpmw.Logging(mux)

	log.Println("server listening on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
