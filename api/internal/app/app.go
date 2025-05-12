package app

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	// "github.com/lotarv/dozens_bot/internal/auth"
	"github.com/lotarv/dozens_bot/internal/config"

	// "github.com/lotarv/dozens_bot/internal/domains/bot"
	"github.com/lotarv/dozens_bot/internal/domains/documents"
	"github.com/lotarv/dozens_bot/internal/domains/members"
	"github.com/lotarv/dozens_bot/internal/domains/notionSync"
	"github.com/lotarv/dozens_bot/internal/domains/user"
	"github.com/lotarv/dozens_bot/internal/storage"
	"github.com/spf13/viper"
)

type controller interface {
	Build()
	Run()
}

type App struct {
	server      *http.Server
	controllers []controller
}

func (app *App) AddController(c controller) {
	app.controllers = append(app.controllers, c)
}

func New() *App {
	config.MustInit(".env")
	app := &App{}
	router := chi.NewMux()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Telegram-Init-Data"},
		ExposedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	router.Use(c.Handler)
	// router.Use(auth.NewAuthMiddleWare())

	server := &http.Server{
		Addr:    "0.0.0.0:" + viper.GetString("port"),
		Handler: router,
	}

	app.server = server

	storage, err := storage.New()
	if err != nil {
		panic(err)
	}

	userController := user.NewUserController(router, storage)
	app.AddController(userController)

	// botController := bot.NewBotController()
	// app.AddController(botController)

	membersController := members.NewMembersController(router, storage)
	app.AddController(membersController)

	notionSyncController := notionSync.NewNotionSyncController(router, storage)
	app.AddController(notionSyncController)

	documentsController := documents.NewDocumentsController(router, storage)
	app.AddController(documentsController)

	return app
}

func (app *App) Init() *App {
	for _, c := range app.controllers {
		c.Build()
	}
	return app
}

func (app *App) Run() {
	slog.Info("Server started add " + app.server.Addr)
	for _, c := range app.controllers {
		go c.Run()
	}
	if err := app.server.ListenAndServe(); err != nil {
		slog.Error(err.Error())
	}
}
