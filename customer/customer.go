package customer

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type Customer struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

func CreateHandler(c *gin.Context) {
	cusReq := Customer{}
	if err := c.ShouldBind(&cusReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO customers (name, email, status) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	if err := stmt.QueryRow(cusReq.Name, cusReq.Email, cusReq.Status).Scan(&cusReq.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cusReq)
}
