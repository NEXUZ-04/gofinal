package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	database "github.com/NEXUZ-04/gofinal/Database"
	customer "github.com/NEXUZ-04/gofinal/Model"
	"github.com/gin-gonic/gin"
)

var db database.DB

func init() {

	// Test hardcode : PASS
	// err := db.Connect("postgres://oeenwpjw:zQqtsCaL5VoY3x-NX8dbzfH7unkJ9Lb0@hansken.db.elephantsql.com:5432/oeenwpjw")

	// Test os.Setenv : PASS
	//os.Setenv("DATABASE_URL", "postgres://oeenwpjw:zQqtsCaL5VoY3x-NX8dbzfH7unkJ9Lb0@hansken.db.elephantsql.com:5432/oeenwpjw")

	fmt.Println(os.Getenv("DATABASE_URL"))
	err := db.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}

	db.Table = "Customers"
	err = db.CreateTB()
	if err != nil {
		log.Fatal("Cannot create table: ", err)
	}
}

func main() {

	defer db.Abort()
	result := gin.Default()

	result.GET("/customers", QueryAllHandler)
	result.GET("/customers/:id", QueryByIdHandler)
	result.POST("/customers", InsertHandler)
	result.PUT("/customers/:id", UpdateHandler)
	result.DELETE("/customers/:id", DeleteHandler)

	result.Run(":2009")
}

func QueryAllHandler(c *gin.Context) {

	pfs, err := db.QueryAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, pfs)
}

func QueryByIdHandler(c *gin.Context) {

	rowId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var pf customer.Profile
	pf, err = db.Query(rowId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, pf)
}

func InsertHandler(c *gin.Context) {

	var pf customer.Profile
	var err error

	if err = c.ShouldBindJSON(&pf); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	pf, err = db.Insert(pf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, pf)
}

func UpdateHandler(c *gin.Context) {

	var pf customer.Profile
	rowId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err = c.ShouldBindJSON(&pf); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	pf.ID = rowId

	pf, err = db.Update(pf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, pf)
}

func DeleteHandler(c *gin.Context) {

	type DeleteResp struct {
		Message string `json:"message"`
	}

	rowId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err = db.Delete(rowId); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, DeleteResp{Message: "customer deleted"})
}
