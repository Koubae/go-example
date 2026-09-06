package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func addRoutes(mux *http.ServeMux, app *Application) {
	mux.Handle("/", http.NotFoundHandler())
	mux.Handle("GET /health", healthcheckHandler(app))

	mux.Handle("GET /accounts", authMiddleware(accountGetHandler(app)))
}

func healthcheckHandler(app *Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := app.container.db.Raw("SELECT 1").Error; err != nil {
			app.logger.Error("healthcheck failed", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			app.logger.Error("healthcheck failed", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}

// export TOKEN=$(uuidgen)
//
//	curl http://localhost:8080/accounts   -H "Authorization: Bearer $TOKEN"
func accountGetHandler(app *Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		accountID := ctx.Value(contextKeyAccountID).(uuid.UUID)
		if accountID == uuid.Nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		response := map[string]string{"accountID": accountID.String()}
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			app.logger.Error("account get failed", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}
