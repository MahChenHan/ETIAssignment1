# ETI Assignment1

## Overview

This Go-based project is a car-pooling platform featuring user authentication, trip creation, and management functionalities. Users can register, log in, create car-pooling trips, join existing trips, and manage their profiles.

## Features

* User registration and authentication
* Trip creation and management
* User profile updates and deletions
* Trip joining and cancellation

## User Authentication
Users can log in with their user IDs.
New users can create an account, specifying whether they are car owners.

## Trip Management
Car owners can create new trips, specifying details such as pickup location, destination, and scheduled start time.
Users can view available trips and join them if there are available seats.
Trip owners can update the scheduled start time.
Trips can be deleted by the owner within a 30-minute window before the scheduled start time.

## Profile Management
Users can update their information, including name, mobile number, and email address.
Car owners provide additional details such as driver's license and car plate number during registration.
Accounts can be deleted, but only if the user has been inactive for at least one year.

## Data Persistence

MySQL is used for the persistent storage of user and trip information.

## Microservices Architecture
### User Service: Manages user authentication, registration, and profile management.
### Trip Service: Handles the creation, retrieval, and management of car-pooling trips.
### Database Service: Manages interactions with the MySQL database, ensuring data consistency across services.
### Authentication Service: Provides token-based authentication for secure communication between microservices.

### Architecture Diagram

            +-----------------------+
            |      User Service      |
            +-----------------------+
                      |
            +-----------------------+
            |      Trip Service      |
            +-----------------------+
                      |
            +-----------------------+
            |   Database Service    |
            +-----------------------+
                      |
            +-----------------------+
            | Authentication Service |
            +-----------------------+

# Setting Up and Running Microservices

## Prerequisites
Go installed on your system
MySQL database

## Installation
Clone the repository from GitHub.
Import the provided SQL script (carpooling.sql) into your MySQL database.

## Running Microservices

### User Service 

Navigate to the user_service directory and run go run console.go.

### Trip Service

Navigate to the trip_service directory and run go run console.go.

### Database Service

Navigate to the database_service directory and run go run console.go.

### Authentication Service 

Navigate to the auth_service directory and run go run console.go.

Note: Ensure that the database connection strings in each microservice's code match your MySQL database settings.

## API Documentation
Detailed API documentation for each microservice is available in their respective directories.

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
