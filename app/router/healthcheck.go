package router

import "net/http"

// HealthCheckHandler ensures the base endpoint is available and can return an HTTP status code 200
func HealthCheckHandler() AppHandler {
	return func(w http.ResponseWriter, r *http.Request) *AppError {
		w.WriteHeader(http.StatusOK)

		msg := "GCP-VISION-API is Healthy :) - Authored by Martin Ombura Jr.\nGitHub: @martinomburajr"
		_, err := w.Write([]byte(msg))
		if err != nil {
			return AppErrorf(http.StatusInternalServerError, "error writing to client", err)
		}
		return nil
	}
}
