package v1

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

func CreateTicketOne(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ticket Ticket
		var sqlStatement string
		if err := c.BindJSON(&ticket); err != nil {
			c.String(http.StatusBadRequest, "Failed to bind request")
			return
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
                time_exit timestamp DEFAULT '0001-01-01',
                vehicle varchar NOT NULL,
                parking_slot varchar NOT NULL
            )`
		if _, err := db.Exec(sqlStatement); err != nil {
			errString := fmt.Sprintf("Error creating database table: %q", err)
			c.String(http.StatusInternalServerError, errString)
			return
		}

		//Insert new ticket into tickets
		ticket.TimeEntry = time.Now()
		sqlStatement = `
			INSERT INTO tickets (vehicle, parking_slot) 
			VALUES ($1, $2) RETURNING id, time_entry`
		if err := db.QueryRow(sqlStatement,
			ticket.Vehicle,
			ticket.ParkingSlot,
		).Scan(&ticket.Id, &ticket.TimeEntry); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error inserting ticket: %q", err))
		}

		c.JSON(http.StatusOK, ticket)
	}
}

func GetTickets(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tickets := make([]Ticket, 0)
		rows, err := db.Query("SELECT id, time_entry, vehicle, parking_slot FROM tickets")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error reading ticks: %q", err))
			return
		}

		defer rows.Close()

		for rows.Next() {
			var ticket Ticket
			if err := rows.Scan(&ticket.Id, &ticket.TimeEntry, &ticket.Vehicle, &ticket.ParkingSlot); err != nil {
				c.String(http.StatusInternalServerError, fmt.Sprintf("Error scanning ticket: %q", err))
				return
			}
			tickets = append(tickets, ticket)
		}

		c.JSON(http.StatusOK, tickets)
	}
}

func GetTicketOne(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if len(id) == 0 {
			c.String(http.StatusBadRequest, "Missing id")
			return
		}

		var ticket Ticket
		err := db.QueryRow(
			"SELECT id, time_entry, vehicle, parking_slot FROM tickets WHERE id = $1", id,
		).Scan(&ticket.Id, &ticket.TimeEntry, &ticket.Vehicle, &ticket.ParkingSlot)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("%q", err))
			return
		}

		c.JSON(http.StatusOK, ticket)
	}
}

func UpdateTicketOne(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Get id request
		id := c.Param("id")
		if len(id) == 0 {
			c.String(http.StatusBadRequest, "Missing id")
			return
		}

		//Parse ticket
		var ticket Ticket
		if err := c.BindJSON(&ticket); err != nil {
			c.String(http.StatusBadRequest, "Failed to bind request")
			return
		}

		//Check nonempty new values
		sqlNewValues := make([]string, 0)
		if len(ticket.ParkingSlot) > 0 {
			sqlNewValues = append(sqlNewValues, fmt.Sprintf("parking_slot='%s'", ticket.ParkingSlot))
		}
		if len(ticket.Vehicle) > 0 {
			sqlNewValues = append(sqlNewValues, fmt.Sprintf("vehicle='%s'", ticket.Vehicle))
		}
		if !ticket.TimeEntry.IsZero() {
			sqlNewValues = append(sqlNewValues, fmt.Sprintf("time_entry='%s', ", ticket.TimeEntry.Format(time.RFC3339Nano)))
		}
		if !ticket.TimeExit.IsZero() {
			sqlNewValues = append(sqlNewValues, fmt.Sprintf("time_exit='%s'", ticket.TimeExit.Format(time.RFC3339Nano)))
		}
		if len(sqlNewValues) == 0 {
			c.String(http.StatusBadRequest, "Missing update values")
			return
		}

		//Update new values into database
		sqlStatement := fmt.Sprintf(`
UPDATE tickets
SET %s
WHERE id = '%s'
RETURNING id, time_entry, time_exit, vehicle, parking_slot`, strings.Join(sqlNewValues, ", "), id)
		fmt.Println(sqlStatement)
		err := db.QueryRow(sqlStatement).Scan(
			&ticket.Id, &ticket.TimeEntry, &ticket.TimeExit, &ticket.Vehicle, &ticket.ParkingSlot,
		)
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, ticket)
	}
}
