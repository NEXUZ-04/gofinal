package main

import (
	"fmt"
	"database/sql"
	"github.com/NEXUZ-04/gofinal/database"
)

var Database DB

func main() {

	err := Database.Connect()
	
}
