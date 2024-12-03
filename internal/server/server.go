package httpserver

import (
	"dist_db/internal/db"
	"dist_db/internal/raft"
	"encoding/json"
	"fmt"
	"net/http"
)

// HTTPServer представляет собой обработчик HTTP-запросов к распределенной БД
type HTTPServer struct {
	replicationManager *raft.Replicator
}

// NewHTTPServer создает новый экземпляр HTTP-сервера
func NewHTTPServer(raftNode *raft.RaftServer) *HTTPServer {
	return &HTTPServer{
		replicationManager: raft.NewReplicator(raftNode),
	}
}

// HandleGet обрабатывает GET-запросы для получения значения по ключу
func (srv *HTTPServer) HandleGet(w http.ResponseWriter, r *http.Request) {
	requestedKey := r.PathValue("key")
	if requestedKey == "" {
		http.Error(w, "Key is not specified in the request", http.StatusBadRequest)
		return
	}

	if storedValue, exists := db.Get(requestedKey); exists {
		w.Write([]byte(fmt.Sprintf("%v", storedValue)))
		return
	}

	http.Error(w, "Record not found", http.StatusNotFound)
}

// HandleCreate обрабатывает POST-запросы для создания новой записи
func (srv *HTTPServer) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	ok, err := srv.replicationManager.ApplyAndReplicate("CREATE", payload.Key, &payload.Value, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error during operation: %v", err), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("%v", ok)))
}

// HandleDelete обрабатывает DELETE-запросы для удаления записи
func (srv *HTTPServer) HandleDelete(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Key string `json:"key"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	ok, err := srv.replicationManager.ApplyAndReplicate("DELETE", payload.Key, nil, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error during deletion: %v", err), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("%v", ok)))
}

// HandleUpdate обрабатывает PUT-запросы для обновления существующей записи
func (srv *HTTPServer) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	ok, err := srv.replicationManager.ApplyAndReplicate("UPDATE", payload.Key, &payload.Value, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error during update: %v", err), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("%v", ok)))
}

// HandleCompareAndSwap обрабатывает PATCH-запросы для атомарного обновления
func (srv *HTTPServer) HandleCompareAndSwap(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Key      string `json:"key"`
		OldValue string `json:"old_value"`
		NewValue string `json:"new_value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	ok, err := srv.replicationManager.ApplyAndReplicate("CAS", payload.Key, &payload.NewValue, &payload.OldValue)
	if err != nil {
		http.Error(w, "Error during CAS", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("%v", ok)))
}
