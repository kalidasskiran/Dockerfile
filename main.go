package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Abc struct {
	Name string `json:"name"`
}

var as []Abc

func ListStructHandler(c *gin.Context) {

	c.JSON(http.StatusOK, as)

}
func InsertStructHandler(c *gin.Context) {
	var a Abc
	c.ShouldBindJSON(&a)
	as = append(as, a)
	c.JSON(http.StatusOK, a)

}

func main() {
	router := gin.Default()
	router.GET("/", ListStructHandler)
	router.POST("/", InsertStructHandler)
	router.Run()
}
