# ETI Assignment1

## Overview

This project is a car-pooling platform implemented in Go, featuring a microservices architecture. The platform enables users to register, publish car-pooling trips, and search for available trips.

## Features

### User Management
Users can register and log in.
Car owners can provide additional information such as driver's license number and car plate number.

### Trip Management
Car owners can publish car-pooling trips with details like pick-up locations, start times, destinations, and available seats.
Passengers can search and enrol in available trips.

### Data Persistence

MySQL is used for the persistent storage of user and trip information.

### Main Application
The main application simulates the car-pooling platform's front end.
It interacts with the secondary microservice to perform necessary actions.

## Microservices
### Main Microservice
Prompts inputs and acts as a 'console'.
Handles user choices such as logging in, enrolling in trips, registering, updating information, deleting accounts, and managing trips.
Drivers can publish carpooling trips, start their trips, or cancel hosted trips.

### Database Microservice
Connects the database and the console.
Is called upon by the main application with the necessary inputs to perform functions such as creating, deleting, and updating users.

## Database Tables
### User
* UserID
* FirstName
* LastName
* MobileNumber
* EmailAddress
* IsCarOwner
* LicenseNumber
* CarPlateNumber
* DateOfCreation

### Trip
* TripID
* PickupLocations
* StartTravelingTime
* DestinationAddress
* MaxPassengers
* EnrolledPassengers
* CarOwnerID
* TripStatus
 
### PassengerTrip
* PassengerTripID
* PassengerID
* PassengerEmail
* DriverID
* DriverEmail
* TripID
* TripCompleted
