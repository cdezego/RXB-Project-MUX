package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Set up the connection information for the posgres database
const (
	host     = "localhost"
	port     = 5555
	user     = "postgres"
	password = "postgres"
	dbname   = "dvdrental"
)

// The Film structure that will be returned via the API
type Film struct {
	Film_ID          int     `json:"film_id"`
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	Release_year     string  `json:"release_year"`
	Language         string  `json:"language"`
	Rental_duration  int     `json:"rental_duration"`
	Rental_rate      float32 `json:"rental_rate"`
	Length           int     `json:"length"`
	Replacement_cost float32 `json:"replacement_cost"`
	Rating           string  `json:"rating"`
	Category         string  `json:"category"`
}

// The welcome function, not much to see here!
func welcomeHandler(c *gin.Context) {

	type KeyValue struct {
		key   string
		value string
	}

	type Params struct {
		items []KeyValue
	}

	data := Params{items: []KeyValue{
		{key: "title", value: "Thingy"},
		{key: "body", value: "Testing 123"},
	}}

	for _, val := range data.items {
		fmt.Println(val.key + ":" + val.value)
	}

	c.IndentedJSON(http.StatusOK, "Welcome to Mockbuster!")
}

// A function to help clean up the code (for building the SQL Where clause)
// The function should be called once for each available search parameter
func addWhereClause(currentString string, pName string, pValue string) string {
	cReturn := currentString

	// Make sure there is a "value" in the search parameters
	if len(pValue) > 0 && pValue != "'" {
		// If the current clause is empty, start it out with WHERE (otherwise add using AND)
		if len(currentString) == 0 {
			cReturn = " WHERE " + pName + strings.ToUpper(pValue)
		} else {
			cReturn = cReturn + " AND " + pName + strings.ToUpper(pValue)
		}
	}

	return cReturn
}

// Write one GET API method that can accept optional parameters (Film_ID, Title, Rating, & Category) for returning a subset of films
// If no parameters are passed, all records are returned and up to all 4 parameters can be used at the same time
// For example, this would utilize all parameters:
//    $ curl -G http://localhost:8080/getFilms?Title='ZORRO%20Ark'\&Category='ComedY'\&Film_ID=1000\&Rating='nc-17'
//
// NOTE: The search is case InsENsitiVE

func getFilms(c *gin.Context) {
	// Connect to the database
	psqlInfo := fmt.Sprintf("host=%s port=%d  dbname=%s user=%s password=%s sslmode=disable", host, port, dbname, user, password)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	// Build the SQL statement
	sqlStatement := "SELECT " +
		"f.film_id," +
		"f.title," +
		"f.description," +
		"f.release_year," +
		"l.name AS language, " +
		"f.rental_duration," +
		"f.rental_rate," +
		"f.length," +
		"f.replacement_cost," +
		"f.rating," +
		"c.name AS category " +
		"FROM category c " +
		"INNER JOIN film_category fc ON c.category_id = fc.category_id " +
		"INNER JOIN film f ON fc.film_id = f.film_id " +
		"INNER JOIN language l ON f.language_id=l.language_id"

	// The Where clause is built separate
	wclause := ""
	wclause = addWhereClause(wclause, "f.film_id=", c.Query("Film_ID"))
	wclause = addWhereClause(wclause, "UPPER(f.title)='", c.Query("Title")+"'")
	wclause = addWhereClause(wclause, "f.rating='", c.Query("Rating")+"'")
	wclause = addWhereClause(wclause, "UPPER(c.Name)='", c.Query("Category")+"'")

	// Then added to complete the statement
	sqlStatement = sqlStatement + wclause // + " ORDER BY film_id"

	fmt.Println("")
	fmt.Println("Running SQL: " + "\n" + sqlStatement)
	fmt.Println("")

	// Run the SQL
	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}

	// Declare the variables needed for the JSON to be returned (for one or more films)
	var fs []Film
	var rate, cost float32
	var Film_ID, duration, length int
	var title, year, rating, description, category, language string

	for rows.Next() {
		rows.Scan(&Film_ID, &title, &description, &year, &language, &duration, &rate, &length, &cost, &rating, &category)
		fs = append(fs, Film{Film_ID, title, description, year, language, duration, rate, length, cost, rating, category})
	}

	db.Close()

	c.IndentedJSON(http.StatusOK, fs)
}

// @title User API documentation
// @version 1.0.0
// @host localhost:8080
// @BasePath /
func main() {
	// Connect to the database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Ping the database to ensure it is responding
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// The database connection succeded
	fmt.Println("Successfully connected to the " + dbname + " database!")

	// Use GIN to start the routes
	router := gin.Default()
	router.GET("/", welcomeHandler)
	router.GET("/getFilms", getFilms)

	http.ListenAndServe(":8080", router)
}
