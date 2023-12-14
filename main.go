package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
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
	DateOfCreation string
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

var currentUser string
var currentName string
var currentType string
var currentID string

func main() {

	// Connect to the MySQL database
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/CarPoolingDB")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Verify the connection
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Simulate user registration and login
	_, err = registerOrLoginUser(db)
	if err != nil {
		log.Fatal(err)
	}

	// if user is a driver,
	if currentType == "Driver" {
		fmt.Println("Welcome " + currentName + " to our carpooling app!")
		fmt.Println("What would you like to do today?")
		fmt.Println("1. Publish a carpooling trip")
		fmt.Println("2. Start a carpooling trip")
		fmt.Println("3. Cancel a carpooling trip")
		fmt.Println("4. Update your profile")
		fmt.Println("5. Delete your account")
		fmt.Println("0. Exit")
		var dChoice int
		fmt.Print("Select an option: ")
		fmt.Scan(&dChoice)

		switch dChoice {
		case 1:
			// Create Trip
			err := createTrip(db, currentUser)
			if err != nil {
				fmt.Print(err)
				return
			}

		case 2:

			// Start a trip, first displaying trips
			err := displayDriverTrips(db)
			if err != nil {
				fmt.Print(err)
				return
			}

			// Choose trip to start
			var tripChoice string
			fmt.Print("Which carpooling trip would you like to start? ")
			fmt.Scan(&tripChoice)

			// Start the trip
			err = startTrip(db, currentUser, tripChoice)
			if err != nil {
				fmt.Print(err)
				return
			}

		case 3:

			// Cancel a trip, first displaying trips
			err := displayDriverTrips(db)
			if err != nil {
				fmt.Print(err)
				return
			}

			// Choosing which trip to cancel
			var tripChoice string
			fmt.Print("Which carpooling trip would you like to cancel? ")
			fmt.Scan(&tripChoice)

			err = cancelTrip(db, currentUser, tripChoice)
			if err != nil {
				fmt.Print(err)
				return
			}

		case 4:

			// Updating User Profile
			user, err := updateProfile(db, currentUser)
			if err != nil {
				fmt.Print(err)
				return
			}

			// Checking user type again
			currentUser = user.EmailAddress
			if user.IsCarOwner {
				currentType = "Driver"
			} else {
				currentType = "Passenger"
			}

			return

		case 5:

			// Deleting user
			err := deleteProfile(db, currentUser)
			if err != nil {
				fmt.Print(err)
				return
			}

		case 0:
			break

		default:

			fmt.Print("Invalid Choice.")

		}

		// Passenger menu
	} else if currentType == "Passenger" {
		fmt.Println("Welcome " + currentName + " to our carpooling app!")
		fmt.Println("What would you like to do today?")
		fmt.Println("1. View all available carpooling trips")
		fmt.Println("2. View previous carpooling trips")
		fmt.Println("3. Enroll in a carpool trip")
		fmt.Println("4. Update your profile")
		fmt.Println("5. Delete your account")
		fmt.Println("0. Exit")
		var pChoice int
		fmt.Print("Select an option: ")
		fmt.Scan(&pChoice)

		switch pChoice {
		case 1:

			// Displaying all available trips
			displayAllTrips(db)
			return

		case 2:

			// Display trip history
			displayPreviousTrips(db)
			return

		case 3:

			// Display all trips, then prompt user to choose which trip
			// they wish to join
			displayAllTrips(db)
			var tripChoice string
			fmt.Print("Which carpooling trip would you like to join? ")
			fmt.Scan(&tripChoice)

			err := enrollTrip(db, currentUser, tripChoice)
			if err != nil {
				fmt.Print(err)
				return
			}

		case 4:

			// Update User Profile
			user, err := updateProfile(db, currentUser)
			if err != nil {
				fmt.Print(err)
				return
			}

			// Checking user type again
			currentUser = user.EmailAddress
			if user.IsCarOwner {
				currentType = "Driver"
			} else {
				currentType = "Passenger"
			}

			return

		case 5:

			// Deleting user
			err := deleteProfile(db, currentUser)
			if err != nil {
				fmt.Print(err)
				return
			}

		case 0:

			break

		default:

			fmt.Print("Invalid Choice.")
		}

	} else {
		return
	}

}

// Register/Login function
func registerOrLoginUser(db *sql.DB) (*User, error) {
	fmt.Println("\nWelcome to our carpooling trip management app! \nFirstly, please:\n1. Register")
	fmt.Println("2. Login\n0. Exit")

	var choice int
	fmt.Print("Select an option: ")
	fmt.Scan(&choice)

	switch choice {
	case 1:
		user, err := registerUser(db)
		if err != nil {
			return nil, err
		}
		fmt.Println("\nRegistration successful!")
		// Displaying user and trip information
		currentUser = user.EmailAddress
		currentName = user.FirstName
		// Checking if it is a driver or passenger
		if user.IsCarOwner {
			currentType = "Driver"
		} else {
			currentType = "Passenger"
		}
		return user, nil
	case 2:
		user, err := loginUser(db)
		if err != nil {
			return nil, err
		}
		fmt.Println("\nLogin successful!")
		// Displaying user and trip information
		currentUser = user.EmailAddress
		currentName = user.FirstName
		// Checking if it is a driver or passenger
		if user.IsCarOwner {
			currentType = "Driver"
		} else {
			currentType = "Passenger"
		}
		return user, nil
	case 0:
		break

	default:
		return nil, fmt.Errorf("Invalid Choice.")
	}

	return nil, nil
}

// Register User function
func registerUser(db *sql.DB) (*User, error) {

	// prompt for all inputs
	var user User
	fmt.Print("First Name: ")
	fmt.Scan(&user.FirstName)
	fmt.Print("Last Name: ")
	fmt.Scan(&user.LastName)
	fmt.Print("Mobile Number: ")
	fmt.Scan(&user.MobileNumber)
	fmt.Print("Email Address: ")
	fmt.Scan(&user.EmailAddress)

	// Check for duplicate email address
	if isDuplicate, err := isDuplicateEmail(db, user.EmailAddress); err != nil {
		return nil, err
	} else if isDuplicate {
		return nil, fmt.Errorf("Email address already exists")
	}

	// Check for duplicate mobile number
	if isDuplicate, err := isDuplicateMobileNumber(db, user.MobileNumber); err != nil {
		return nil, err
	} else if isDuplicate {
		return nil, fmt.Errorf("Mobile number already exists")
	}

	// Additional registration fields for car owners
	fmt.Print("Are you a car owner? (yes/no): ")
	var isCarOwner string
	fmt.Scan(&isCarOwner)
	if isCarOwner == "yes" {
		user.IsCarOwner = true
		fmt.Print("License Number: ")
		fmt.Scan(&user.LicenseNumber)
		fmt.Print("Car Plate Number: ")
		fmt.Scan(&user.CarPlateNumber)
	}

	postBody, _ := json.Marshal(user)
	resBody := bytes.NewBuffer(postBody)

	// Insert user into the Users table
	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPost, "http://localhost:8000/users/", resBody); err == nil { //Runs POST method using user api link from second Go File
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 { //Status code indicating success
				return &user, nil // Returns user
			}
		}
	}

	return &user, nil
}

// Login function
func loginUser(db *sql.DB) (*User, error) {
	var email string

	fmt.Print("Email Address: ")
	fmt.Scan(&email)

	// Query user from the Users table
	row := db.QueryRow("SELECT FirstName, LastName, MobileNumber, EmailAddress, IsCarOwner, LicenseNumber, CarPlateNumber, DateOfCreation FROM User WHERE EmailAddress = ?", email)

	var user User
	err := row.Scan(&user.FirstName, &user.LastName, &user.MobileNumber, &user.EmailAddress, &user.IsCarOwner, &user.LicenseNumber, &user.CarPlateNumber, &user.DateOfCreation)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Function to check for duplicate email address
func isDuplicateEmail(db *sql.DB, email string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM User WHERE EmailAddress = ?", email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Function to check for duplicate mobile number
func isDuplicateMobileNumber(db *sql.DB, mobileNumber string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM User WHERE MobileNumber = ?", mobileNumber).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Function to display all trips that are available
func displayAllTrips(db *sql.DB) {

	// Filtering through to get only available trips
	rows, err := db.Query("SELECT * FROM Trip WHERE TripStatus = true AND EnrolledPassengers < MaxPassengers")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var trips []Trip

	for rows.Next() {
		var trip Trip
		err := rows.Scan(
			&trip.TripID,
			&trip.PickupLocations,
			&trip.StartTravelingTime,
			&trip.DestinationAddress,
			&trip.MaxPassengers,
			&trip.EnrolledPassengers,
			&trip.CarOwnerID,
			&trip.TripStatus,
		)
		if err != nil {
			log.Fatal(err)
		}

		trips = append(trips, trip)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Display all trips
	fmt.Println("\nAll Available Trips")
	for _, trip := range trips {
		fmt.Println("-------------------------------")
		fmt.Printf("TripID: %s\n", trip.TripID)
		fmt.Printf("Pickup Location: %s\n", trip.PickupLocations)
		fmt.Printf("Start Travel Time: %s\n", trip.StartTravelingTime)
		fmt.Printf("Destination: %s\n", trip.DestinationAddress)
		fmt.Printf("Max Passengers: %d\n", trip.MaxPassengers)
		fmt.Printf("Enrolled Passengers: %d\n", trip.EnrolledPassengers)
		fmt.Printf("Car Owner UserID: %s\n", trip.CarOwnerID)
	}
}

// Display trip history based on email
func displayPreviousTrips(db *sql.DB) {
	rows, err := db.Query(`
	SELECT PickupLocations, StartTravelingTime, DestinationAddress, MaxPassengers, EnrolledPassengers, CarOwnerID FROM Trip INNER JOIN PassengerTrip ON PassengerTrip.TripID = Trip.TripID WHERE TripCompleted = true AND PassengerEmail = ? ORDER BY CAST(Trip.StartTravelingTime as datetime) DESC`, currentUser)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var trips []Trip

	for rows.Next() {
		var trip Trip
		err := rows.Scan(
			&trip.PickupLocations,
			&trip.StartTravelingTime,
			&trip.DestinationAddress,
			&trip.MaxPassengers,
			&trip.EnrolledPassengers,
			&trip.CarOwnerID,
		)
		if err != nil {
			log.Fatal(err)
		}

		trips = append(trips, trip)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Display trips
	fmt.Println("\nYour Previous Trips")
	for _, trip := range trips {
		fmt.Println("Your Previous Trips-------------------------------")
		fmt.Printf("Pickup Location: %s\n", trip.PickupLocations)
		fmt.Printf("Start Travel Time: %s\n", trip.StartTravelingTime)
		fmt.Printf("Destination: %s\n", trip.DestinationAddress)
		fmt.Printf("Max Passengers: %d\n", trip.MaxPassengers)
		fmt.Printf("Enrolled Passengers: %d\n", trip.EnrolledPassengers)
		fmt.Printf("Car Owner UserID: %s\n", trip.CarOwnerID)
	}
}

// Display driver trips based on email
func displayDriverTrips(db *sql.DB) error {

	// Filtering through, getting incomplete trips, filtered by email and timestamp, ensuring it is 30 minutes around scheduled time.
	rows, err := db.Query(`
	SELECT Trip.TripID, PickupLocations, StartTravelingTime, DestinationAddress, MaxPassengers, EnrolledPassengers, CarOwnerID FROM Trip INNER JOIN PassengerTrip ON PassengerTrip.TripID = Trip.TripID WHERE TripCompleted = false AND DriverEmail = ? AND TIMESTAMPDIFF(MINUTE, NOW(), Trip.StartTravelingTime) >= 30 ORDER BY CAST(Trip.StartTravelingTime as datetime) DESC`, currentUser)
	if err != nil {
		fmt.Print("You have no trips that are able to be cancelled or started.")
		log.Fatal(err)
		return err
	}

	defer rows.Close()

	var trips []Trip

	for rows.Next() {
		var trip Trip
		err := rows.Scan(
			&trip.TripID,
			&trip.PickupLocations,
			&trip.StartTravelingTime,
			&trip.DestinationAddress,
			&trip.MaxPassengers,
			&trip.EnrolledPassengers,
			&trip.CarOwnerID,
		)
		if err != nil {
			log.Fatal(err)
			return (err)
		}

		trips = append(trips, trip)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Display trips
	fmt.Println("\nYour Published Trips")
	if len(trips) == 0 {
		fmt.Print("You have no available trips to start/cancel.\n")
		return nil
	}

	for _, trip := range trips {
		fmt.Println("-------------------------------")
		fmt.Printf("Trip ID: %s\n", trip.TripID)
		fmt.Printf("Pickup Location: %s\n", trip.PickupLocations)
		fmt.Printf("Start Travel Time: %s\n", trip.StartTravelingTime)
		fmt.Printf("Destination: %s\n", trip.DestinationAddress)
		fmt.Printf("Max Passengers: %d\n", trip.MaxPassengers)
		fmt.Printf("Enrolled Passengers: %d\n", trip.EnrolledPassengers)
		fmt.Printf("Car Owner UserID: %s\n", trip.CarOwnerID)
	}

	return nil
}

// Update Profile function
func updateProfile(db *sql.DB, email string) (*User, error) {

	var UserID string
	var user User

	// Finding user
	result := db.QueryRow("select UserID, FirstName, LastName, MobileNumber, EmailAddress, IsCarOwner, LicenseNumber, CarPlateNumber from User where EmailAddress = ?", currentUser) // Retrieves user records where UserID matches current session UserID
	err := result.Scan(&UserID, user.FirstName, user.LastName, user.MobileNumber, user.EmailAddress, user.IsCarOwner, user.LicenseNumber, user.CarPlateNumber)
	if err == sql.ErrNoRows {
		fmt.Println("No Rows")
	}

	// Prompt for updates
	fmt.Print("Update First Name: ")
	fmt.Scan(&user.FirstName)

	fmt.Print("Update Last Name: ")
	fmt.Scan(&user.LastName)

	fmt.Print("Update Mobile Number: ")
	fmt.Scan(&user.MobileNumber)

	fmt.Print("Update Email Address: ")
	fmt.Scan(&user.EmailAddress)

	// Additional registration fields for car owners
	fmt.Print("Are you a car owner? (yes/no): ")
	var isCarOwner string
	fmt.Scan(&isCarOwner)
	if isCarOwner == "yes" {
		user.IsCarOwner = true
		fmt.Print("License Number: ")
		fmt.Scan(&user.LicenseNumber)
		fmt.Print("Car Plate Number: ")
		fmt.Scan(&user.CarPlateNumber)
	}

	fmt.Print(&user.MobileNumber)

	postBody, _ := json.Marshal(user)
	resBody := bytes.NewBuffer(postBody)

	// Insert user into the Users table
	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPut, "http://localhost:8000/users/"+UserID, resBody); err == nil { // Runs Put method using user api link from second Go File
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 {
				return &user, nil // Returns user
			}
		}
	}

	return &user, nil
}

// Delete Profile function
func deleteProfile(db *sql.DB, email string) error {

	var user User
	var UserID string

	// Check if the account is over 1 year old
	err := db.QueryRow("SELECT * FROM User WHERE EmailAddress = ? AND TIMESTAMPDIFF(YEAR, DateOfCreation, NOW()) >= 1;", email).Scan(&UserID, user.FirstName, user.LastName, user.MobileNumber, user.EmailAddress, user.IsCarOwner, user.LicenseNumber, user.CarPlateNumber, user.DateOfCreation)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Your account must be at least 1 year old to delete.")
		} else {
			return err
		}
	}

	// Perform the deletion from second microservice
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8000/users/%d", UserID), nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		fmt.Println("Account deleted successfully")
		currentName = ""
		currentType = ""
		currentUser = ""
		return nil
	}

	return fmt.Errorf("Failed to delete account")
}

// Join Trip function
func enrollTrip(db *sql.DB, email string, tripID string) error {

	var isExists bool
	var passengerID string
	var driverEmail string
	var driverID string
	var maxPassengers int
	var enrolledPassengers int
	var pTrip PassengerTrip

	// Checking if trip exists
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM Trip WHERE TripID = ?)", tripID).Scan(&isExists)
	if err != nil {
		log.Fatal(err)
	}

	if !isExists {
		return fmt.Errorf("Trip ID does not exist")
	}

	// Getting DriverID
	err = db.QueryRow("SELECT CarOwnerID FROM Trip WHERE TripID = ?", tripID).Scan(&driverID)
	if err != nil {
		log.Fatal(err)
	}

	// Getting PassengerID
	err = db.QueryRow("SELECT UserID FROM User WHERE EmailAddress = ?", email).Scan(&passengerID)
	if err != nil {
		log.Fatal(err)
	}

	err = db.QueryRow("SELECT MaxPassengers, EnrolledPassengers FROM Trip WHERE TripID = ?", tripID).Scan(&maxPassengers, &enrolledPassengers)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the maximum number of passengers has been reached
	if enrolledPassengers >= maxPassengers {
		fmt.Print("The trip is full. Cannot enroll more passengers.")
		return fmt.Errorf("The trip is full. Cannot enroll more passengers.")
	}

	// Getting DriverEmail
	err = db.QueryRow("SELECT EmailAddress FROM User WHERE UserID = ? ", driverID).Scan(&driverEmail)
	if err != nil {
		log.Fatal(err)
	}

	// Initialized pTrip Trip variable
	pTrip.PassengerEmail = email
	pTrip.TripCompleted = false // Assuming TripStatus is a boolean field
	pTrip.DriverEmail = driverEmail
	pTrip.DriverID = driverID
	pTrip.PassengerID = passengerID
	pTrip.TripID = tripID

	// Check if the user is already enrolled in the trip
	var isEnrolled bool
	err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM PassengerTrip WHERE PassengerID = ? AND TripID = ?)", passengerID, tripID).Scan(&isEnrolled)
	if err != nil {
		log.Fatal(err)
	}

	if isEnrolled {
		fmt.Printf("User is already enrolled in the trip.")
		return fmt.Errorf("User is already enrolled in the trip")
	}

	// Update the EnrolledPassengers count for the trip
	_, err = db.Exec("UPDATE Trip SET EnrolledPassengers = EnrolledPassengers + 1 WHERE TripID = ?", tripID)
	if err != nil {
		log.Fatal(err)
	}

	postBody, _ := json.Marshal(pTrip)
	resBody := bytes.NewBuffer(postBody)

	// Insert user into the Users table
	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPost, "http://localhost:8000/ptrip/", resBody); err == nil { //Runs POST method using user api link from second Go File
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 { //Status code indicating success
				return nil // Returns null
			}
		}
	}

	fmt.Printf("UserID %s enrolled in the trip with ID %s\n", passengerID, tripID)
	return nil
}

// Create Trip function
func createTrip(db *sql.DB, email string) error {

	var trip Trip

	// Retrieve car owner ID using the provided email
	err := db.QueryRow("SELECT UserID FROM User WHERE EmailAddress = ?", email).Scan(&trip.CarOwnerID)
	if err != nil {
		return fmt.Errorf("Failed to retrieve car owner ID: %v", err)
	}

	// Getting necessary inputs
	fmt.Print("Pickup Locations: ")
	fmt.Scan(&trip.PickupLocations)

	fmt.Print("Start Traveling Time (YYYY-MM-DD|HH:MM): ")
	//var timeStr string
	fmt.Scan(&trip.StartTravelingTime)
	//trip.StartTravelingTime, err = time.Parse("2006-01-02 15:04", timeStr)

	if err != nil {
		return fmt.Errorf("Failed to parse time: %v", err)
	}

	fmt.Print("Destination Address: ")
	fmt.Scan(&trip.DestinationAddress)

	fmt.Print("Max Passengers: ")
	fmt.Scan(&trip.MaxPassengers)

	postBody, _ := json.Marshal(trip)
	resBody := bytes.NewBuffer(postBody)

	// Insert user into the Users table
	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPost, "http://localhost:8000/trip/", resBody); err == nil { //Runs POST method using user api link from second Go File
		if _, err := client.Do(req); err == nil {
		}
	}

	fmt.Printf("Trip created successfully\n")
	return nil
}

// Starting Trip function
func startTrip(db *sql.DB, email string, tripID string) error {

	var isExists bool
	var isDriver bool
	var trip Trip

	// Checking if trip exists
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM Trip WHERE TripID = ?)", tripID).Scan(&isExists)
	if err != nil {
		log.Fatal(err)
	}

	if !isExists {
		return fmt.Errorf("Trip does not exist")
	}

	// Checking if user is the driver
	err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM PassengerTrip WHERE TripID = ? AND DriverEmail = ?)", tripID, currentUser).Scan(&isDriver)
	if err != nil {
		log.Fatal(err)
	}

	if !isDriver {
		return fmt.Errorf("You are not the driver for this trip.")
	}

	// Checking if it is within 30 minutes
	err = db.QueryRow("SELECT * FROM Trip WHERE TripID = ? AND TIMESTAMPDIFF(MINUTE, NOW(), StartTravelingTime) >= 30;", tripID).Scan(&trip.TripID, &trip.PickupLocations, &trip.StartTravelingTime, &trip.DestinationAddress, &trip.MaxPassengers, &trip.EnrolledPassengers, &trip.CarOwnerID, &trip.TripStatus)
	if err != nil {
		print("It is not within the 30 minutes timeframe.")
		log.Fatal(err)
		return err
	}

	// Updating TripStatus
	trip.TripStatus = false

	postBody, _ := json.Marshal(trip)
	resBody := bytes.NewBuffer(postBody)

	// Updating trip
	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPut, "http://localhost:8000/trip/"+trip.TripID, resBody); err == nil { // Runs Put method using user api link from second Go File
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 { //Status code indicating success
				fmt.Print("Trip officially started")
				return nil
			}
		}
	}

	return nil

}

// function to cancel trip
func cancelTrip(db *sql.DB, email string, tripID string) error {

	// Check if it is within the 30 minute timeframe
	var trip Trip

	err := db.QueryRow("SELECT * FROM Trip WHERE TripID = ? AND TIMESTAMPDIFF(MINUTE, NOW(), StartTravelingTime) >= 30;", tripID).Scan(&trip.TripID, &trip.PickupLocations, &trip.StartTravelingTime, &trip.DestinationAddress, &trip.MaxPassengers, &trip.EnrolledPassengers, &trip.CarOwnerID, &trip.TripStatus)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("The trip is not within the 30 minute timeframe.")
			return err

		} else {
			return err
		}
	}

	// Perform the deletion from the second microservice
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8000/trip/"+tripID, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		fmt.Println("Trip cancelled successfully")
		return nil
	}

	return fmt.Errorf("Failed to cancel trip")

}

// Simple function to display user information
func displayUserInfo(user User) {
	fmt.Printf("User Information:\n"+
		"Name: %s %s\n"+
		"Mobile Number: %s\n"+
		"Email Address: %s\n"+
		"Is Car Owner: %t\n"+
		"License Number: %s\n"+
		"Car Plate Number: %s\n",
		user.FirstName, user.LastName, user.MobileNumber, user.EmailAddress, user.IsCarOwner, user.LicenseNumber, user.CarPlateNumber)
}
