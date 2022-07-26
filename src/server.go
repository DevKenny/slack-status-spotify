package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"os"

	"github.com/DevKenny/slack-spotify/src/handlers"
	"github.com/DevKenny/slack-spotify/src/repositories"
	"github.com/DevKenny/slack-spotify/src/repositories/db_entities"
	"github.com/DevKenny/slack-spotify/src/services"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/robfig/cron/v3"
	"github.com/zmb3/spotify"
	"gorm.io/driver/postgres"
)

func main() {
	// Get environment variables
	newRelicAppName := os.Getenv("SPOTIFY_SLACK_APP_NEW_RELIC_APP_NAME")
	newRelicLicense := os.Getenv("SPOTIFY_SLACK_APP_NEW_RELIC_LICENSE")
	databaseURL := os.Getenv("SPOTIFY_SLACK_APP_DATABASE_URL")
	slackAuthURL := os.Getenv("SPOTIFY_SLACK_APP_SLACK_AUTH_URL")
	spotifyRedirectURL := os.Getenv("SPOTIFY_SLACK_APP_SPOTIFY_REDIRECT_URL")
	slackClientID := os.Getenv("SPOTIFY_SLACK_APP_SLACK_CLIENT_ID")
	slackClientSecret := os.Getenv("SPOTIFY_SLACK_APP_SLACK_CLIENT_SECRET")
	port := os.Getenv("PORT")

	// Setup New Relic
	newRelicApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(newRelicAppName),
		newrelic.ConfigLicense(newRelicLicense),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	// Setup connection to the database
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&db_entities.User{})

	// Creating Spotify Authenticator
	spotifyAuthenticator := spotify.NewAuthenticator(spotifyRedirectURL, spotify.ScopeUserReadCurrentlyPlaying)

	// Creating app layers (repositories, services, handlers)
	repositories := repositories.NewRepository(db)
	services := services.NewServices(repositories, spotifyAuthenticator)
	handlers := handlers.NewHandlers(services, spotifyAuthenticator, stateGenerator(), slackClientID, slackClientSecret, slackAuthURL)

	// Setup cronjob for updating status
	c := cron.New(cron.WithSeconds())
	c.AddFunc("@every 10s", func() { services.ChangeUserStatus(context.Background()) })
	c.Start()

	// Add handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", handlers.SpotifyCallbackHandler)
	mux.HandleFunc("/slackAuth", handlers.SlackCallbackHandler)
	mux.HandleFunc(newrelic.WrapHandleFunc(newRelicApp, "/users", handlers.HealthHandler))
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)

	http.ListenAndServe(":"+port, mux)
}

func stateGenerator() string {
	b := make([]byte, 4)
	i, err := rand.Read(b)
	if err != nil {
		return "0"
	}
	state := fmt.Sprintf("%x", b)
	println("Spotify State random => ", i)
	println("Spotify State => ", state)
	return state
}
