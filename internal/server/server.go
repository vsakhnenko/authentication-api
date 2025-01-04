package server

import (
	"authentication/internal/config"
	impl2 "authentication/internal/repositories/impl"
	"authentication/internal/services"
	"authentication/internal/services/impl"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"authentication/internal/database"
)

type Server struct {
	port        int
	db          database.Service
	authService services.AuthService
}

func NewServer() *http.Server {
	port, err := strconv.Atoi(os.Getenv("BACKEND_PORT"))
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	config.LoadEnvVariables()
	config.InitRedis()
	config.SMTPConnect()
	db := database.New()

	//repositories
	repository := impl2.NewUserRepository(db.DB())

	//services
	authSer := impl.NewAuthService(repository)

	serverInstance := &Server{
		port:        port,
		db:          db,
		authService: authSer,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", serverInstance.port),
		Handler:      serverInstance.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
