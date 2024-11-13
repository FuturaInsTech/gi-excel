package main

import (
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/FuturaInsTech/GoExcel/initializers"

	"github.com/FuturaInsTech/GoExcel/routes"
)

func init() {
	initializers.LoadEnvVariables()
	username := os.Getenv("USER_NAME")
	password := os.Getenv("PASSWORD")
	initializers.ConnectToDb(username, password, "localhost", "3306")
	initializers.SyncDatabase()
}

func main() {

	r := gin.Default()

	// Configure CORS to dynamically allow any origin with credentials
	config := cors.Config{
		AllowCredentials: true, // Allow credentials (cookies, authorization headers, etc.)
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// Dynamically set AllowOrigins to reflect the origin of the incoming request
		AllowOriginFunc: func(origin string) bool {
			return true // Allow all origins
		},
	}

	r.Use(cors.New(config))

	// Serve frontend static files from the dist directory (React build)
	r.Static("/static", "./static/dist")
	r.LoadHTMLFiles("./static/dist/index.html")

	// Serve API routes
	v1 := r.Group("/api/v1")

	{
		routes.Authroutes(v1)
		routes.Basicroutes(v1)
		routes.ExcelRoutes(v1)
	}

	// Default route to serve the React app (all unmatched routes)
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Fetch the port from environment variables
	goport := os.Getenv("PORT")
	if goport == "" {
		goport = "3002" // Default goport if PORT is not specified in the .env file
	}

	// Start the server on the specified goport
	r.Run(":" + goport)
}
