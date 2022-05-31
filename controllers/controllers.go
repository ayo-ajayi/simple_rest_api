package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/twinj/uuid"
)

type Choice struct {
	ID   string `json:"id"`
	Gone bool   `json:"gone"`
	Come bool   `json:"come"`
}

var db *sql.DB
var err error

func newId() string {
	return uuid.NewV4().String()
}

var DBinit = func() {
	if err = godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening DB: %s", err.Error())
	}
	db.Ping()
	if err = db.Ping(); err != nil {
		log.Fatalf("Could not ping: %s", err.Error())
	}
	if _, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS choice(
			id VARCHAR(36),
			go BOOL,
			come BOOL
		)`); err != nil {
		log.Fatalf("Could not execute table creation: %s", err.Error())
	}
	log.Println("Connected to DB successfully")
}

var PostChoice = func(c *gin.Context) {
	newChoice := Choice{
		ID:   newId(),
		Gone: false,
		Come: false,
	}
	err := c.BindJSON(&newChoice)
	if err != nil {
		log.Printf("invalid input: %v", err)
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	log.Printf("{%s, %v, %v}", newChoice.ID, newChoice.Gone, newChoice.Come)
	validate := validator.New()
	validationErr := validate.Struct(&newChoice)
	if validationErr != nil {
		c.JSON(400, gin.H{"error": validationErr.Error()})
		return
	}
	_, err = db.Exec(`insert into choice(id, go, come) values($1, $2, $3)`, newChoice.ID, newChoice.Gone, newChoice.Come)
	if err != nil {
		log.Fatalf("DB code could not run with success: %s", err.Error())
	}

	m := "successfully created the record"
	log.Println(m)
	c.JSON(200, gin.H{
		"message": m,
	})
}

var GetChoice = func(c *gin.Context) {
	rows, err := db.Query(`SELECT id, go, come FROM choice`)
	switch err {
	case sql.ErrNoRows:
		defer rows.Close()
		e := "no rows records found in choice table to read"
		log.Println(e)
		c.JSON(400, gin.H{"error": e})
	case nil:
		defer rows.Close()
		a := make([]Choice, 0)
		var rowsReadErr bool
		for rows.Next() {
			var id string
			var gone, come bool
			err = rows.Scan(&id, &gone, &come)
			if err != nil {
				log.Printf("error occurred while reading the database rows: %v", err)
				rowsReadErr = true
				break
			}
			a = append(a, Choice{id, gone, come})
		}

		if rowsReadErr {
			log.Println("we are not able to fetch few records")
		}
		log.Printf("we are able to fetch choices")
		c.JSON(200, a)

	default:
		defer rows.Close()
		e := "some internal database server error"
		log.Println(e)
		c.JSON(500, gin.H{"error": e})

	}
}

var CheckID = func(c *gin.Context) {
	id := c.Param("id")
	var gone, come bool
	row := db.QueryRow(`SELECT id, go, come FROM choice WHERE id= $1 LIMIT 1`, id)
	err = row.Scan(&id, &gone, &come)
	switch err {
	case sql.ErrNoRows:
		e := fmt.Sprintf("row with id %v not found in choice table", id)
		log.Println(e)
		c.JSON(500, gin.H{"error": e})
	case nil:
		log.Println("we are able to fetch the choice")
		a := Choice{id, gone, come}
		c.Set("values", a)
		c.Next()
	default:
		e := "some internal database server error"
		log.Println(e)
		c.JSON(500, gin.H{"error": e})
	}
}

var GetChoiceByID = func(c *gin.Context) {
	res := c.MustGet("values").(Choice)
	log.Println(res)
	c.JSON(200, &res)
}

var UpdateChoice = func(c *gin.Context) {
	res := c.MustGet("values").(Choice).ID
	updateChoice := Choice{}
	err := c.BindJSON(&updateChoice)
	if err != nil {
		log.Printf("invalid input: %v", err)
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	log.Printf("{%s, %v, %v}", updateChoice.ID, updateChoice.Gone, updateChoice.Come)
	updateChoice.ID = res
	_, err = db.Exec(`UPDATE choice SET go=$1, come=$2 WHERE id=$3`, updateChoice.Gone, updateChoice.Come, res)
	if err != nil {
		log.Fatalf("DB code could not run with success: %s", err.Error())
	}

	m := "successfully updated the record"
	log.Println(m)
	c.JSON(200, gin.H{"message": m,
		"record": updateChoice})
}

var DeleteChoice = func (c *gin.Context){
	res := c.MustGet("values").(Choice).ID
	_, err = db.Exec(`delete from choice where id = $1`, res)
	if err != nil {
		log.Fatalf("DB code could not run with success: %s", err.Error())
	}

	m := "successfully deleted the record"
	log.Println(m)
	c.JSON(200, gin.H{
		"message": m,
	})
}