package main

import (
	"net/http"

	"github.com/Cingozr/go-rest-api/internal/comment"
	"github.com/Cingozr/go-rest-api/internal/database"
	transportHTTP "github.com/Cingozr/go-rest-api/internal/transport/http"
	log "github.com/sirupsen/logrus"
)

//App - the striuct which contains thins like ponters
//to database connection
type App struct {
	Name    string
	Version string
}

//Run - sets up Our Application
func (app *App) Run() error {
	log.SetFormatter(&log.JSONFormatter{})
	log.WithFields(
		log.Fields{
			"AppName": app.Name,
			"Version": app.Version,
		}).Info("Setting up application")

	db, err := database.NewDatabase()
	if err != nil {
		return err
	}

	err = database.MigrateDb(db)
	if err != nil {
		return err
	}

	commentService := comment.NewService(db)

	handler := transportHTTP.NewHandler(commentService)
	handler.SetupRoutes()

	if err := http.ListenAndServe(":8080", handler.Router); err != nil {
		log.Error("Failed to set up server")
		return err
	}

	return nil
}

// func init() {
// 	os.Setenv("DB_USERNAME", "postgres")
// 	os.Setenv("DB_PASSWORD", "postgres")
// 	os.Setenv("DB_HOST", "localhost")
// 	os.Setenv("DB_PORT", "5432")
// 	os.Setenv("DB_TABLE", "postgres")
// }

func main() {
	app := App{
		Name:    "Comments Service",
		Version: "1.0.0",
	}

	if err := app.Run(); err != nil {
		log.Error("Error starting up our REST API")
		log.Fatal(err)
	}
}
