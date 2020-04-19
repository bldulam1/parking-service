package v1

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func CreateTicket(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		//vehicle, parkingSlot := c.Query("vehicle"), c.Query("parkingSlot")
		var ticket Ticket
		var sqlStatement string
		if err := c.BindJSON(&ticket); err != nil {
			c.String(http.StatusBadRequest, "Failed to bind request")
		}
		//Validate request
		if len(ticket.ParkingSlot) == 0 {
			c.String(http.StatusBadRequest, "Missing Parking Slot")
			return
		}
		if len(ticket.Vehicle) == 0 {
			c.String(http.StatusBadRequest, "Missing Vehicle Information")
			return
		}

		//Create table if not exists
		sqlStatement = `
            CREATE TABLE IF NOT EXISTS tickets (
                id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
                time_entry timestamp DEFAULT now(),
                time_exit timestamp,
                vehicle varchar NOT NULL,
                parkingSlot varchar NOT NULL
            )`
		if _, err := db.Exec(sqlStatement); err != nil {
			errString := fmt.Sprintf("Error creating database table: %q", err)
			c.String(http.StatusInternalServerError, errString)
			return
		}

		//Insert new ticket into tickets
		ticket.TimeEntry = time.Now()
		sqlStatement = `
			INSERT INTO tickets (vehicle, parkingSlot, time_entry) 
			VALUES ($1, $2, $3) RETURNING id`
		if err := db.QueryRow(sqlStatement,
			ticket.Vehicle,
			ticket.ParkingSlot,
			ticket.TimeEntry.Format(time.RFC3339Nano),
		).Scan(&ticket.Id); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error inserting ticket: %q", err))
		}
		fmt.Println(ticket.Id, ticket.TimeEntry)

		c.JSON(http.StatusOK, string(ticket.JSON()))
	}
}

func GetTickets(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sqlStatement string
		sqlStatement = "SELECT * FROM tickets"

		rows, err := db.Query(sqlStatement)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error reading ticks: %q", err))
			return
		}

		defer rows.Close()

		for rows.Next() {
			var id uuid.UUID
			if err := rows.Scan(&id); err != nil {
				c.String(http.StatusInternalServerError, fmt.Sprintf("Error scanning ticket: %q", err))
				return
			}
			fmt.Println(id)
		}

		c.String(http.StatusOK, fmt.Sprintf("Read from DB: \n"))
	}
}
