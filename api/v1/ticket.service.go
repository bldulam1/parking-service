package v1

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateTicket(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ticket := NewTicket(c.Query("vehicle"), c.Query("parkingSlot"))
		c.JSON(http.StatusOK, string(ticket.JSON()))

		if _, err := db.Exec(`
			CREATE TABLE IF NOT EXISTS tickets (
				id uuid DEFAULT uuid_generate_v4(), 
				timeEntry timestamp, 
				timeExit timestamp, 
				vehicle varchar, 
				parkingSlot varchar
			)
		`); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error creating database table: %q", err))
			return
		}

		if _, err := db.Exec(fmt.Sprintf(`
			INSERT INTO tickets (timeEntry, vehicle, parkingSlot) 
				VALUES (now(),%s,%s)`, ticket.Vehicle, ticket.ParkingSlot)); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error incrementing tick: %q", err))
			return
		}
		//
		//rows, err := db.Query("SELECT tick FROM ticks")
		//if err != nil {
		//	c.String(http.StatusInternalServerError,
		//		fmt.Sprintf("Error reading ticks: %q", err))
		//	return
		//}
		//
		//defer rows.Close()
		//for rows.Next() {
		//	var tick time.TimeEntry
		//	if err := rows.Scan(&tick); err != nil {
		//		c.String(http.StatusInternalServerError,
		//			fmt.Sprintf("Error scanning ticks: %q", err))
		//		return
		//	}
		//	c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", tick.String()))
		//}
	}
}

func GetTickets(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
