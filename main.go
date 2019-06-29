package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"

	"github.com/spfrank01/finalexam/customer"
)

func main() {
	r := gin.Default()
	r.Use(addMiddleware)
	r.POST("/customers", customer.CreateHandler)
	r.GET("/customers/:id", customer.GetByIdHandler)
	r.GET("/customers", customer.GetHandler)
	r.PUT("/customers/:id", customer.UpdateByIDHandler)
	r.DELETE("/customers/:id", customer.DeleteByIdHandler)
	r.Run(":2019")
}

func addMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != "token2019" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
		c.Abort()
		return
	}
	c.Next()
}
