package main

import (
	"encoding/json"
	"fmt"
	
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Person struct {
	gorm.Model

	Name  string
	Email string `gorm:"typevarchar(100);unique_index"`
	RollNo string
}

var db *gorm.DB
var err error


var (
	person=&Person{
		Model:  gorm.Model{},
		Name:   "Debayan",
		Email:  "debayanbiswas1111@gmail.com",
		RollNo: "21PH10012",
	}
)

func main(){
	// Loading enviroment variables
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbname := os.Getenv("NAME")
	dbpassword := os.Getenv("PASSWORD")

	// Database connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbname, dbpassword, dbPort)

	// Openning connection to database
	db, err = gorm.Open(dialect, dbURI)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connected to database successfully")
	}

	// Close the databse connection when the main function closes
	defer db.Close()

	// Make migrations to the database if they haven't been made already
	db.AutoMigrate(&Person{})

	db.Create(&person)

	/*----------- API routes ------------*/
	router := mux.NewRouter()
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/person/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/create/person", CreatePerson).Methods("POST")
	router.HandleFunc("/delete/person/{id}", DeletePerson).Methods("DELETE")

	http.ListenAndServe(":8080", router)
	

}
/*-------- API Controllers --------*/

/*----- People ------*/
func GetPeople(w http.ResponseWriter, r *http.Request) {
	var people []Person

	db.Find(&people)

	json.NewEncoder(w).Encode(&people)
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	

	db.First(&person, params["id"])
	

	json.NewEncoder(w).Encode(&person)
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var person Person
	json.NewDecoder(r.Body).Decode(&person)

	createdPerson := db.Create(&person)
	err = createdPerson.Error
	if err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(&createdPerson)
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person Person

	db.First(&person, params["id"])
	db.Delete(&person)

	json.NewEncoder(w).Encode(&person)
}