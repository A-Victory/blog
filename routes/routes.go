package routes

import (
	"net/http"

	"github.com/A-Victory/blog/auth"
	"github.com/A-Victory/blog/database/conn"
	"github.com/A-Victory/blog/handlers"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/cors"
)

type ServerConfig struct {
	DB *conn.DB
	VA *auth.Validation
}

func NewServer(config ServerConfig) *chi.Mux {
	router := chi.NewRouter()

	// Middlewares
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowCredentials: false,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Authorization"},
		Debug:            true,
	}).Handler)
	router.Use(setJSONContentType)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	handler := handlers.NewHttpHandler(&handlers.Config{
		Database:  config.DB,
		Validator: config.VA,
	})

	router.Get("/health", healthCheck)

	router.Route("/api", func(r chi.Router) {

		registerRoutes(r, handler)

		authRouter := r.With(auth.Verify)

		authRouter.Get("/users/profile", handler.Profile)

		postRoutes(authRouter, handler)

		commentRoutes(authRouter, handler)

	})

	return router
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	render.Data(w, r, []byte("Ok"))
}

func setJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func registerRoutes(r chi.Router, httpHandler *handlers.HttpHandler) {
	r.Route("/users", func(router chi.Router) {
		router.Post("/register", httpHandler.CreateUser)
		router.Post("/login", httpHandler.Login)
	})
}

func postRoutes(r chi.Router, httpHandler *handlers.HttpHandler) {
	r.Route("/posts", func(router chi.Router) {
		router.Get("/", httpHandler.Post)
		router.Get("/{id}", httpHandler.Post)
		router.Put("/{id}", httpHandler.Post)
		router.Delete("/{id}", httpHandler.Post)
		router.Post("/", httpHandler.Post)
	})
}

func commentRoutes(r chi.Router, httpHandler *handlers.HttpHandler) {
	r.Route("/posts/{postId}/comments", func(router chi.Router) {
		router.Route("/", func(route chi.Router) {
			route.Get("/", httpHandler.Comment)
			route.Post("/", httpHandler.Comment)
		})
	})
	r.Route("/comments", func(route chi.Router) {
		route.Put("/{id}", httpHandler.Comment)
		route.Delete("/{id}", httpHandler.Comment)
	})
}
