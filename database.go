package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// User represents the user data structure.
type User struct {
	FirstName      string
	LastName       string
	MobileNumber   string
	EmailAddress   string
	IsCarOwner     bool
	LicenseNumber  string
	CarPlateNumber string
}

// Trip represents the car-pooling trip data structure.
type Trip struct {
	TripID             string
	PickupLocations    string
	StartTravelingTime string
	DestinationAddress string
	MaxPassengers      int
	EnrolledPassengers int
	CarOwnerID         string
	TripStatus         bool
}

// PassengerTrip represents the passenger trip data structure.
type PassengerTrip struct {
	PassengerID    string
	PassengerEmail string
	DriverID       string
	DriverEmail    string
	TripID         string
	TripCompleted  bool
}

var (
	db  *sql.DB
	err error
)

func main() {
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/CarPoolingDB")

	if err != nil {
		panic(err.Error())
	}

	router := mux.NewRouter()
	router.HandleFunc("/users/", user).Methods("POST")
	router.HandleFunc("/users/{userid}", user).Methods("DELETE", "PUT")
	router.HandleFunc("/trip/", trip).Methods("POST")
	router.HandleFunc("/trip/{tripid}", trip).Methods("DELETE", "PUT")
	router.HandleFunc("/ptrip/", passengertrip).Methods("POST")
	router.HandleFunc("/ptrip/{passengertripid}", passengertrip).Methods("PUT")
	fmt.Println("Listening at port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

// user Function
func user(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// if it is creating a user,
	if r.Method == "POST" {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data User
			fmt.Println(string(body))
			if err := json.Unmarshal(body, &data); err == nil {
				fmt.Println(data)
				insertUser(data)
				w.WriteHeader(http.StatusAccepted)
			} else {
				fmt.Println(err)
			}
		}

		// if it is updating a user,
	} else if r.Method == "PUT" {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data User

			if err := json.Unmarshal(body, &data); err == nil {
				fmt.Println(data)
				updateUser(params["userid"], data)
				w.WriteHeader(http.StatusAccepted)
			} else {
				fmt.Println(err)
			}
		}

		// if it is deleting a user,
	} else if r.Method == "DELETE" {
		fmt.Fprintf(w, params["userid"]+" Deleted")
		deleteUser(params["userid"])
	} else {
	}
}

// Trip function
func trip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// if it is creating a trip
	if r.Method == "POST" {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data Trip
			fmt.Println(string(body))
			if err := json.Unmarshal(body, &data); err == nil {
				fmt.Println(data)
				insertTrip(data)
				w.WriteHeader(http.StatusAccepted)
			} else {
				fmt.Println(err)
			}
		}
		// if it is updating a trip
	} else if r.Method == "PUT" {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data Trip

			if err := json.Unmarshal(body, &data); err == nil {
				fmt.Println(data)
				updateTrip(params["tripid"], data)
				w.WriteHeader(http.StatusAccepted)
			} else {
				fmt.Println(err)
			}
		}
		// if it is deleting a trip
	} else if r.Method == "DELETE" {
		fmt.Fprintf(w, params["tripid"]+" Deleted")
		deleteTrip(params["tripid"])

	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// PassengerTrip function
func passengertrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// if it is creating a PassengerTrip
	if r.Method == "POST" {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data PassengerTrip
			fmt.Println(string(body))
			if err := json.Unmarshal(body, &data); err == nil {
				fmt.Println(data)
				insertPassengerTrip(data)
				w.WriteHeader(http.StatusAccepted)
			} else {
				fmt.Println(err)
			}
		}

		// if it is updating a PassengerTrip
	} else if r.Method == "PUT" {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data PassengerTrip

			if err := json.Unmarshal(body, &data); err == nil {
				fmt.Println(params["passengertripid"], data)
				w.WriteHeader(http.StatusAccepted)
			} else {
				fmt.Println(err)
			}
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// function to delete user
func deleteUser(id string) (int64, error) {
	result, err := db.Exec("delete from User where UserID=?", id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// function to create user
func insertUser(u User) {
	_, err := db.Exec("insert into User (FirstName, LastName, MobileNumber, EmailAddress, IsCarOwner, LicenseNumber, CarPlateNumber, DateOfCreation) values(?, ?, ?, ?, ?, ?, ?, NOW())", u.FirstName, u.LastName, u.MobileNumber, u.EmailAddress, u.IsCarOwner, u.LicenseNumber, u.CarPlateNumber)
	if err != nil {
		panic(err.Error())
	}
}

// function to update user
func updateUser(id string, u User) {
	_, err := db.Exec("update User set FirstName = ?, LastName = ?, MobileNumber = ?, EmailAddress = ?, IsCarOwner = ?, LicenseNumber = ?, CarPlateNumber = ? where UserID=?", u.FirstName, u.LastName, u.MobileNumber, u.EmailAddress, u.IsCarOwner, u.LicenseNumber, u.CarPlateNumber, id)
	if err != nil {
		panic(err.Error())
	}
}

// function to create trip
func insertTrip(t Trip) {
	_, err := db.Exec("insert into Trip (PickUpLocations, StartTravelingTime, DestinationAddress, MaxPassengers, EnrolledPassengers, CarOwnerID) values(?,?,?,?,?,?)", t.PickupLocations, t.StartTravelingTime, t.DestinationAddress, t.MaxPassengers, t.EnrolledPassengers, t.CarOwnerID)
	if err != nil {
		panic(err.Error())
	}
}

// function to update trip
func updateTrip(id string, t Trip) {
	_, err := db.Exec("update Trip set PickUpLocations = ?, StartTravelingTime = ?, DestinationAddress = ?, MaxPassengers = ?, EnrolledPassengers = ?, CarOwnerID = ?, TripStatus = ? where TripID = ?", t.PickupLocations, t.StartTravelingTime, t.DestinationAddress, t.MaxPassengers, t.EnrolledPassengers, t.CarOwnerID, t.TripStatus, id)
	if err != nil {
		panic(err.Error())
	}
}

// function to delete trip
func deleteTrip(id string) (int64, error) {
	result, err := db.Exec("delete from Trip where TripID = ?", id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// function to create passengerTrip
func insertPassengerTrip(p PassengerTrip) {
	_, err := db.Exec("insert into PassengerTrip (PassengerID, PassengerEmail, DriverID, DriverEmail, TripID, TripCompleted) values (?,?,?,?,?,?)", p.PassengerID, p.PassengerEmail, p.DriverID, p.DriverEmail, p.TripID, p.TripCompleted)
	if err != nil {
		panic(err.Error())
	}
}

// function to update passengerTrip
func updatePassengerTrip(id string, p PassengerTrip, status string) {
	_, err := db.Exec("update PassengerTrip set TripStatus = ? WHERE CarPoolID = ?", status, id)
	if err != nil {
		panic(err.Error())
	}
}
