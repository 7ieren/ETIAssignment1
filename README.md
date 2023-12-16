# ETIAssignment1
Overview
This is a car-pooling platform implemented in Go with a microservices architecture. The platform allows users to register, publish car-pooling trips, and search for available trips.

### Architecture
This program has two microservices, and the main.go (Main microservice) acts as the 'front-end' of the program, and calls upon the database.go (Secondary microservice) for insertion, updating, or deletion of records in the database itself. 

### Set up for microservices
To set up the microservice, you are to first run database.go, and ensure that it is listening to port 8000. Afterwards, run main.go to use the platform itself.

### Features
#### User Management:

Users can register and log in.
Car owners can provide additional information such as driver's license number and car plate number.

#### Trip Management:

Car owners can publish car-pooling trips with details like pick-up locations, start times, destinations, and available seats.
Passengers can search and enroll in available trips.

### Data Persistence:
The application uses MySQL for persistent storage of user and trip information.

### Main Application:

A main application is provided to simulate the car-pooling platform front end.
The main app interacts with the secondary microservice to perform necessary actions.

### Microservices
#### Main Microservice:
Prompts inputs, and acts as a 'console', getting user's choices such as logging in, and enrolling in trips.
Users can register, log in, update their information, delete their accounts, and enroll in or view all available trips.
Drivers can publish carpooling trips, start their own trips, or cancel the trips they are hosting, 
on top of updating and deleting their accounts.

#### Database Microservice:
Acts as a connector between the database and the console is called upon by the main application,
with the necessary inputs to perform functions such as but not limited to creating, deleting, and updating users.

### Database
Consists of three tables, User, Trip, and PassengerTrip
### User has the columns: 
  - UserID
  - FirstName
  - LastName 
  - MobileNumber
  - EmailAddress
  - IsCarOwner
  - LicenseNumber
  - CarPlateNumber
  - DateOfCreation
    
### Trip has the columns:
  - TripID
  - PickupLocations
  - StartTravelingTime
  - DestinationAddress
  - MaxPassengers
  - EnrolledPassengers
  - CarOwnerID
  - TripStatus
    
### PassengerTrip has the columns:
  - PassengerTripID
  - PassengerID
  - PassengerEmail
  - DriverID
  - DriverEmail
  - TripID
  - TripCompleted
