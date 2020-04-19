package main

import (
	"bytes"
	"database/sql"
	v1 "github.com/bldulam1/parking-service/api/v1"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
)

func repeatHandler(r int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var buffer bytes.Buffer
		for i := 0; i < r; i++ {
			buffer.WriteString("Hello from Go!\n")
		}
		c.String(http.StatusOK, buffer.String())
	}
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	tStr := os.Getenv("REPEAT")
	repeat, err := strconv.Atoi(tStr)
	if err != nil {
		log.Printf("Error converting $REPEAT to an int: %q - Using default\n", err)
		repeat = 5
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

	router.GET("/mark", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world2")
	})

	router.GET("/repeat", repeatHandler(repeat))

	router.POST("/api/v1/tickets", v1.CreateTicket(db))
	router.GET("/api/v1/tickets", v1.GetTickets(db))
	router.GET("/api/v1/ticket/:id", v1.GetTicket(db))
	router.PUT("/api/v1/ticket/:id", v1.GetTicket(db))

	router.Run(":" + port)
}
