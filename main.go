package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"openchat/api"
	"openchat/lib"
	"openchat/model"
	"openchat/utils/rcontext"
)

func init() {
	lib.InstallDB()
	lib.InstallS3FromEnv()
	lib.InstallDeepseekFromEnv()

	if err := lib.DB.AutoMigrate(
		&model.User{},
		&model.RefreshToken{},
		&model.Conversation{},
		&model.Message{},
		&model.Image{},
	); err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("started...")
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           3600,
	}).Handler)
	apiRouter := chi.NewRouter()
	apiRouter.Post("/signup", api.Auth.Signup)
	apiRouter.Post("/signin", api.Auth.Signin)
	apiRouter.Post("/refresh_token", api.Auth.RefreshToken)
	apiRouter.Group(func(r chi.Router) {
		r.Use(rcontext.WithUserId)
		r.Get("/account", api.Auth.Account)
		r.Get("/conversations", api.Conversation.ListConversations)
		r.Post("/conversations", api.Conversation.CreateConversation)
		r.Get("/conversations/{id}", api.Conversation.GetConversation)
		r.Delete("/conversations/{id}", api.Conversation.DeleteConversation)
		r.Get("/conversations/{id}/messages", api.Message.ListMessages)
		r.Post("/chat", api.Chat.ChatCompletion)
		r.Post("/summarize", api.Chat.Summarize)
		r.Get("/images", api.Image.ListImages)
		r.Delete("/images/{id}", api.Image.DeleteImage)
		r.Post("/cf/images", api.Image.GenerateImageFromCf)
	})
	r.Mount("/api", apiRouter)
	log.Fatal(http.ListenAndServe(":8080", r))
}
