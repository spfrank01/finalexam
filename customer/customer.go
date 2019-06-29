package customer

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Customer struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

func GetHandler(c *gin.Context) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT id, name, email, status FROM customers;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rows, err := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var cusRes []Customer
	for rows.Next() {
		cusTmp := Customer{}

		if err := rows.Scan(&cusTmp.ID, &cusTmp.Name, &cusTmp.Email, &cusTmp.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cusRes = append(cusRes, cusTmp)
	}
	c.JSON(http.StatusOK, cusRes)
}

func GetByIdHandler(c *gin.Context) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT id, name, email, status FROM customers WHERE id=$1;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	row := stmt.QueryRow(c.Param("id"))
	cusRes := Customer{}
	if err := row.Scan(&cusRes.ID, &cusRes.Name, &cusRes.Email, &cusRes.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
	}

	c.JSON(http.StatusOK, cusRes)
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

//UPDATE todos SET title=$2, status=$3 WHERE id=$1`, id, title, status)
func UpdateByIDHandler(c *gin.Context) {
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
	stmt, err := db.Prepare("UPDATE customers SET name=$2, email=$3, status=$4 WHERE id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	if err := stmt.QueryRow(cusReq.ID, cusReq.Name, cusReq.Email, cusReq.Status).Scan(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cusReq)
}

func DeleteByIdHandler(c *gin.Context) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM customers WHERE id=$1;")
	if err != nil {
		fmt.Println("err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	if val, err := stmt.Query(c.Param("id")); err != nil {
		fmt.Println("val", val)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
}
