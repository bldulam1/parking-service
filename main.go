package main

import (
	"database/sql"
	v1 "github.com/bldulam1/parking-service/api/v1"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	psqlInfo := os.Getenv("DATABASE_URL")
	if len(psqlInfo) == 0 {
		psqlInfo = "host=localhost port=5432 user=postgres password=tha3nohk! dbname=parking_service sslmode=disable"
	}

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.POST("/api/v1/ticket", v1.CreateTicketOne(db))
	router.PUT("/api/v1/ticket/:id", v1.UpdateTicketOne(db))
	router.GET("/api/v1/ticket/:id", v1.GetTicketOne(db))
	router.GET("/api/v1/tickets", v1.GetTickets(db))

	router.Run(":" + port)
}
