package app

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"spy_cat/internal/app/endpoint"
	"spy_cat/internal/app/middleware"
	"spy_cat/internal/app/repository"
	"spy_cat/internal/app/routes"
	"spy_cat/internal/app/service"
	"time"
)

var DB *gorm.DB

type App struct {
	e               *echo.Echo
	catEndpoint     *endpoint.CatEndpoint
	missionEndpoint *endpoint.MissionEndpoint
	targetEndpoint  *endpoint.TargetEndpoint
}

func New() (*App, error) {
	a := &App{}

	a.e = echo.New()

	middleware.RegisterMiddlewares(a.e)

	err := ConnectDatabase()
	if err != nil {
		return nil, err
	}

	catRepo := repository.NewCatRepository(DB)
	missionRepo := repository.NewMissionRepository(DB)
	targetRepo := repository.NewTargetRepository(DB)

	catService := service.NewCatService(catRepo)
	missionService := service.NewMissionService(missionRepo, catRepo, targetRepo)
	targetService := service.NewTargetService(targetRepo)

	a.catEndpoint = endpoint.NewCatEndpoint(catService)
	a.missionEndpoint = endpoint.NewMissionEndpoint(missionService)
	a.targetEndpoint = endpoint.NewTargetEndpoint(targetService)

	routes.RegisterRoutes(a.e, a.catEndpoint, a.missionEndpoint, a.targetEndpoint)

	return a, nil
}

func (a *App) Run() {
	fmt.Println("Server running")
	err := a.e.Start(":8087")
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectDatabase() error {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	sqlDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	sqlDB, err := sql.Open("postgres", sqlDSN)
	if err != nil {
		return fmt.Errorf("could not open SQL database connection: %v", err)
	}
	defer sqlDB.Close()

	err = waitForDatabase(sqlDB, 20, 5*time.Second)
	if err != nil {
		return fmt.Errorf("database not ready: %v", err)
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("could not open GORM database connection: %v", err)
	}

	return nil
}

func waitForDatabase(db *sql.DB, retries int, delay time.Duration) error {
	for i := 0; i < retries; i++ {
		err := db.Ping()
		if err == nil {
			return nil
		}
		log.Printf("Waiting for database to be ready (%d/%d)...", i+1, retries)
		time.Sleep(delay)
	}
	return fmt.Errorf("database not ready after %d retries", retries)
}
