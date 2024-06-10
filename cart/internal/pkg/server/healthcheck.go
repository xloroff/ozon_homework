package server

import "net/http"

// Healthcheck проверяем работу сервера.
func Healthcheck(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
}
