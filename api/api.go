package api

import (
	"billy/handlers"
	"billy/middleware"
	"billy/utils"
	"log"
	"net/http"
)

type Api struct {
	handler *handlers.Handler
	mux     *http.ServeMux
}

func NewApi(handler handlers.Handler) *Api {
	return &Api{handler: &handler, mux: http.NewServeMux()}
}

func (api *Api) SetupRoutes() {
	api.mux.HandleFunc("POST /login", api.handler.LoginHandler)
	api.mux.HandleFunc("POST /logout", api.handler.LogoutHandler)
	api.mux.HandleFunc("POST /refresh", api.handler.RefreshHandler)

	api.mux.HandleFunc("POST /users", api.handler.CreateUserHandler)
	api.mux.HandleFunc("GET /users", api.handler.GetUsersHandler)
	api.mux.HandleFunc("GET /users/{id}", api.handler.GetUserByIDHandler)

	api.mux.HandleFunc("GET /usage", middleware.Auth(api.handler.GetUsageHandler))

	api.mux.Handle("/api/", http.StripPrefix("/api", api.mux))
	utils.LogRoutes()
}

func (api *Api) Run() error {
	log.Println("Starting server on :8080...")
	return http.ListenAndServe(":8080", api.mux)
}
