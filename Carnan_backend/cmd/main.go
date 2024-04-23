package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/KITTTPOB-bank/Carnan_backend/pkg/routes"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/handlers"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	routes.SetupRouter(r)
	http.Handle("/", r)
	http.ListenAndServe(getPort(),
		handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			// handlers.AllowedOrigins([]string{"http://localhost:3000"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}),
			handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Accept", "AllowedOrigins", "AllowedHeaders", "AllowedMethods", "Content-Length", "Accept-Encoding", "X-Requested-With"}),
		)(r))

}
func getPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8010"
		fmt.Println("No Port In Heroku" + port)
	}
	return ":" + port
}
