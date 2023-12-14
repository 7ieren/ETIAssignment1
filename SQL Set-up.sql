-- Create the CarPoolingDB database
CREATE DATABASE IF NOT EXISTS CarPoolingDB;
USE CarPoolingDB;

-- Create the Users table
CREATE TABLE IF NOT EXISTS User (
    UserID INT AUTO_INCREMENT PRIMARY KEY,
    FirstName VARCHAR(50) NOT NULL,
    LastName VARCHAR(50) NOT NULL,
    MobileNumber VARCHAR(20) NOT NULL,
    EmailAddress VARCHAR(100) NOT NULL UNIQUE,
    IsCarOwner BOOLEAN NOT NULL,
    LicenseNumber VARCHAR(20),
    CarPlateNumber VARCHAR(20),
    DateOfCreation DATETIME NOT NULL
);

-- Create the Trips table
CREATE TABLE IF NOT EXISTS Trip (
    TripID INT AUTO_INCREMENT PRIMARY KEY,
    PickupLocations TEXT NOT NULL,
    StartTravelingTime DATETIME NOT NULL,
    DestinationAddress VARCHAR(100) NOT NULL,
    MaxPassengers INT NOT NULL,
    EnrolledPassengers INT NOT NULL,
    CarOwnerID VARCHAR(100),
    TripStatus BOOL DEFAULT TRUE,
    FOREIGN KEY (CarOwnerID) REFERENCES User(UserID)
);

-- Add the PassengerTrip table
CREATE TABLE IF NOT EXISTS PassengerTrip (
	PassengerTripID INT AUTO_INCREMENT PRIMARY KEY,
    PassengerID INT,
    PassengerEmail VARCHAR(100) NOT NULL,
    DriverID INT,
    DriverEmail VARCHAR(100),
    TripID INT,
    TripCompleted bool DEFAULT FALSE,
    FOREIGN KEY (PassengerID) REFERENCES User(UserID),
    FOREIGN KEY (DriverID) REFERENCES User(UserID)
);