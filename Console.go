package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// User represents the structure of a user in the system
type User struct {
	ID             int       `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	MobileNumber   string    `json:"mobile_number"`
	EmailAddress   string    `json:"email_address"`
	DriverLicense  string    `json:"driver_license,omitempty"`
	CarPlateNumber string    `json:"car_plate_number,omitempty"`
	CarOwner       string    `json:"car_owner"`
	AccountCreated time.Time `json:"account_created"`
	LastUpdated    time.Time `json:"last_updated"`
	AccountDeleted time.Time `json:"account_deleted,omitempty"`
}

// Trip represents the structure of a car-pooling trip
type Trip struct {
	ID                 int       `json:"id"`
	OwnerID            int       `json:"owner_id"`
	PickupLocation     string    `json:"pickup_location"`
	AltPickupLocation  string    `json:"alt_pickup_location,omitempty"`
	Destination        string    `json:"destination"`
	ScheduledStartTime time.Time `json:"scheduled_start_time"`
	MaxPassengers      int       `json:"max_passengers"`
	AvailablePassenger int       `json:"available_passenger"`
}

var db *sql.DB
var err error

func main() {
	// Initialize the database connection
	db, err = sql.Open("mysql", "root:24A887332y04@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var currentUserID int

	// Main application loop
	for {
		// Display main menu if the user is not logged in, otherwise show logged-in menu
		if currentUserID == 0 {
			printMainMenu()
		} else {
			printLoggedInMenu(currentUserID)
		}

		var choice int
		fmt.Print("Enter your choice: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			// Log in or view profile based on user login status
			if currentUserID == 0 {
				loginUser(&currentUserID)
			} else {
				profileMenu(currentUserID)
			}
		case 2:
			// Create user or manage trips based on user login status
			if currentUserID == 0 {
				createUser()
			} else {
				tripMenu(currentUserID)
			}
		case 3:
			// Exit the program if the user is logged in
			if currentUserID != 0 {
				fmt.Println("Exiting program.")
				os.Exit(0)
			}
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

// Function to display the main menu
func printMainMenu() {
	fmt.Println("=== Main Menu ===")
	fmt.Println("1. Log In User")
	fmt.Println("2. Create User")
	fmt.Println("3. Exit")
}

// Function to display the logged-in menu
func printLoggedInMenu(currentUserID int) {
	for {
		fmt.Println("=== Logged In Menu ===")
		fmt.Println("1. Profile")
		fmt.Println("2. Trip")
		fmt.Println("3. Exit")

		var choice int
		fmt.Print("Enter your choice: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			// View and manage user profile
			profileMenu(currentUserID) // Pass currentUserID to profileMenu
		case 2:
			// View and manage trips
			tripMenu(currentUserID) // Pass currentUserID to tripMenu
		case 3:
			fmt.Println("Logging out.")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

// Function to authenticate and log in a user
func loginUser(currentUserID *int) {
	var userID int
	fmt.Print("Enter User ID: ")
	fmt.Scan(&userID)

	// Perform authentication by checking if the user ID exists in the database
	var storedUserID int
	err := db.QueryRow("SELECT id FROM users WHERE id=?", userID).Scan(&storedUserID)

	if err == sql.ErrNoRows {
		fmt.Println("Login failed. User ID not found.")
		return
	} else if err != nil {
		fmt.Println("Error during login:", err)
		return
	}

	fmt.Println("Login successful. Welcome, user ID:", storedUserID)
	*currentUserID = storedUserID
}

// Function to create a new user
func createUser() {
	var newUser User
	fmt.Print("Enter First Name: ")
	fmt.Scan(&newUser.FirstName)
	fmt.Print("Enter Last Name: ")
	fmt.Scan(&newUser.LastName)
	fmt.Print("Enter Mobile Number: ")
	fmt.Scan(&newUser.MobileNumber)
	fmt.Print("Enter Email Address: ")
	fmt.Scan(&newUser.EmailAddress)
	fmt.Print("Are you a car owner? (yes/no): ")
	fmt.Scan(&newUser.CarOwner)
	newUser.CarOwner = strings.ToLower(newUser.CarOwner) // Convert to lowercase for case-insensitivity

	// Check if the input is valid
	if newUser.CarOwner != "yes" && newUser.CarOwner != "no" {
		fmt.Println("Invalid input for car owner. Please enter 'yes' or 'no'.")
		return
	}

	if newUser.CarOwner == "yes" {
		// Additional input for car owners
		fmt.Print("Enter Driver's License Number: ")
		fmt.Scan(&newUser.DriverLicense)
		fmt.Print("Enter Car Plate Number: ")
		fmt.Scan(&newUser.CarPlateNumber)
	}

	// Set timestamps and insert user into the database
	newUser.AccountCreated = time.Now()
	newUser.LastUpdated = newUser.AccountCreated

	result, err := db.Exec("INSERT INTO users (first_name, last_name, mobile_number, email_address, driver_license, car_plate_number, car_owner, account_created, last_updated) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		newUser.FirstName, newUser.LastName, newUser.MobileNumber, newUser.EmailAddress, newUser.DriverLicense, newUser.CarPlateNumber, newUser.CarOwner, newUser.AccountCreated, newUser.LastUpdated)

	if err != nil {
		fmt.Println("Error creating user:", err)
		return
	}

	newUserID, _ := result.LastInsertId()
	newUser.ID = int(newUserID)

	fmt.Println("User created successfully:")
	printUser(newUser)

	// Automatically log in the newly created user
	fmt.Println("Logging in as user ID:", newUser.ID)
	printLoggedInMenu(newUser.ID)
}

// Function to display the profile menu
func profileMenu(currentUserID int) {
	for {
		fmt.Println("=== Profile Menu ===")
		fmt.Println("1. Update User")
		fmt.Println("2. Delete User")
		fmt.Println("3. Back to Main Menu")

		var choice int
		fmt.Print("Enter your choice: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			// Update user details
			updateUser(currentUserID) // Pass currentUserID
		case 2:
			// Delete user account
			deleteUser(currentUserID) // Pass currentUserID
		case 3:
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

// Function to update user details
func updateUser(currentUserID int) {
	// Use the currentUserID instead of prompting for a new ID
	userID := currentUserID

	// similar to createUser, update user details
	var updatedUser User
	fmt.Print("Enter Updated First Name: ")
	fmt.Scan(&updatedUser.FirstName)
	fmt.Print("Enter Updated Last Name: ")
	fmt.Scan(&updatedUser.LastName)
	fmt.Print("Enter Updated Mobile Number: ")
	fmt.Scan(&updatedUser.MobileNumber)
	fmt.Print("Enter Updated Email Address: ")
	fmt.Scan(&updatedUser.EmailAddress)
	fmt.Print("Are you a car owner? (yes/no): ")
	fmt.Scan(&updatedUser.CarOwner)
	updatedUser.CarOwner = strings.ToLower(updatedUser.CarOwner) // Convert to lowercase for case-insensitivity

	// Check if the input is valid
	if updatedUser.CarOwner != "yes" && updatedUser.CarOwner != "no" {
		fmt.Println("Invalid input for car owner. Please enter 'yes' or 'no'.")
		return
	}

	if updatedUser.CarOwner == "yes" {
		// Additional input for car owners
		fmt.Print("Enter Updated Driver's License Number: ")
		fmt.Scan(&updatedUser.DriverLicense)
		fmt.Print("Enter Updated Car Plate Number: ")
		fmt.Scan(&updatedUser.CarPlateNumber)
	}

	updatedUser.LastUpdated = time.Now()

	_, err := db.Exec("UPDATE users SET first_name=?, last_name=?, mobile_number=?, email_address=?, driver_license=?, car_plate_number=?, car_owner=?, last_updated=? WHERE id=?",
		updatedUser.FirstName, updatedUser.LastName, updatedUser.MobileNumber, updatedUser.EmailAddress, updatedUser.DriverLicense, updatedUser.CarPlateNumber, updatedUser.CarOwner, updatedUser.LastUpdated, userID)

	if err != nil {
		fmt.Println("Error updating user:", err)
		return
	}

	fmt.Println("User updated successfully:")
	printUser(updatedUser)
}

// Function to delete a user account
func deleteUser(currentUserID int) {
	// Use the currentUserID instead of prompting for a new ID
	userID := currentUserID

	// Check if the user with the specified ID exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id=?)", userID).Scan(&exists)
	if err != nil {
		fmt.Println("Error checking user existence:", err)
		return
	}

	if !exists {
		fmt.Println("User with ID", userID, "does not exist.")
		return
	}

	// Get the last updated time as a string
	var lastUpdatedStr string
	err = db.QueryRow("SELECT last_updated FROM users WHERE id=?", userID).Scan(&lastUpdatedStr)
	if err != nil {
		fmt.Println("Error getting last updated time:", err)
		return
	}

	// Parse the string into a time.Time variable
	lastUpdated, err := time.Parse("2006-01-02 15:04:05", lastUpdatedStr)
	if err != nil {
		fmt.Println("Error parsing last updated time:", err)
		return
	}

	oneYearAgo := time.Now().AddDate(-1, 0, 0)
	if lastUpdated.After(oneYearAgo) {
		fmt.Println("User has not been inactive for 1 year. Deletion not allowed.")
		return
	}

	// Continue with the deletion by setting AccountDeleted field
	_, err = db.Exec("UPDATE users SET account_deleted=? WHERE id=?", time.Now(), userID)
	if err != nil {
		fmt.Println("Error deleting user:", err)
		return
	}

	// Respond with success message
	fmt.Println("User deleted successfully.")
}

// Function to display the trip menu
func tripMenu(currentUserID int) {
	fmt.Println("=== Trip Menu ===")
	fmt.Println("1. Create Trip")
	fmt.Println("2. Get Trip")
	fmt.Println("3. Update Start Trip")
	fmt.Println("4. Delete Trip")
	fmt.Println("5. Check My Trips")
	fmt.Println("6. Back to Main Menu")

	var choice int
	fmt.Print("Enter your choice: ")
	fmt.Scan(&choice)

	switch choice {
	case 1:
		// Create a new trip
		createTrip(currentUserID)
	case 2:
		// Display available trips
		getTrip(currentUserID)
	case 3:
		// Update start time of a trip
		updateStartTime(currentUserID)
	case 4:
		// Delete a trip
		deleteTrip(currentUserID)
	case 5:
		// Display trips joined by the user
		checkMyTrips(currentUserID)
	case 6:
		// Return after handling the choice
		return
	default:
		fmt.Println("Invalid choice. Please try again.")
	}

	// Return after handling the choice
}

// Function to check if the user is the owner or part of the trip
func isUserOwnerOrParticipant(userID, tripID int) (bool, error) {
	// Check if the user is the owner or participant in the trip
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM trips WHERE id=? AND (owner_id=? OR EXISTS(SELECT 1 FROM user_trips WHERE user_id=? AND trip_id=?)))", tripID, userID, userID, tripID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Function to retrieve all trips a user has taken before in reverse chronological order
func checkMyTrips(currentUserID int) {
	// Retrieve trips joined by the user in reverse chronological order
	rows, err := db.Query("SELECT t.id, t.owner_id, t.pickup_location, t.alt_pickup_location, t.destination, t.scheduled_start_time, t.max_passengers, t.available_passenger FROM trips t JOIN user_trips ut ON t.id = ut.trip_id WHERE ut.user_id = ? ORDER BY t.scheduled_start_time DESC", currentUserID)
	if err != nil {
		fmt.Println("Error checking trips:", err)
		return
	}
	defer rows.Close()

	fmt.Println("Trips You've Joined (Reverse Chronological Order):")
	for rows.Next() {
		var trip Trip
		var scheduledStartTimeStr string

		err := rows.Scan(&trip.ID, &trip.OwnerID, &trip.PickupLocation, &trip.AltPickupLocation,
			&trip.Destination, &scheduledStartTimeStr, &trip.MaxPassengers, &trip.AvailablePassenger,
		)
		if err != nil {
			fmt.Println("Error scanning trip:", err)
			return
		}

		// Parse the scheduled_start_time string into a time.Time variable
		scheduledStartTime, err := time.Parse("2006-01-02 15:04:05", scheduledStartTimeStr)
		if err != nil {
			fmt.Println("Error parsing start time:", err)
			return
		}

		// Assign the parsed time to the trip structure
		trip.ScheduledStartTime = scheduledStartTime

		printTrip(trip)
	}
}

// Function to check if the user is a car owner
func isUserCarOwner(userID int) (bool, error) {
	var carOwnerStr string
	// Retrieve the car_owner field for the user
	err := db.QueryRow("SELECT car_owner FROM users WHERE id=?", userID).Scan(&carOwnerStr)
	if err != nil {
		return false, err
	}

	// Check if the car_owner field is "yes"
	return carOwnerStr == "yes", nil
}

// Function to create a new trip
func isUserTripOwner(userID, tripID int) (bool, error) {
	var ownerID int
	// Retrieve the owner_id for the given trip
	err := db.QueryRow("SELECT owner_id FROM trips WHERE id=?", tripID).Scan(&ownerID)
	if err != nil {
		return false, err
	}
	return userID == ownerID, nil
}

// Function to create a new trip
func createTrip(currentUserID int) {
	var newTrip Trip
	newTrip.OwnerID = currentUserID

	// Check if the user exists before checking if they are a car owner
	var userExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id=?)", newTrip.OwnerID).Scan(&userExists)
	if err != nil {
		fmt.Println("Error checking if user exists:", err)
		return
	}

	if !userExists {
		fmt.Println("User with ID", newTrip.OwnerID, "does not exist.")
		return
	}

	// Check if the owner is a car owner
	isCarOwner, err := isUserCarOwner(newTrip.OwnerID)
	if err != nil {
		fmt.Println("Error checking if user is a car owner:", err)
		return
	}

	if !isCarOwner {
		fmt.Println("Only car owners can publish car-pooling trips.")
		return
	}

	fmt.Print("Enter Pickup Location: ")
	fmt.Scan(&newTrip.PickupLocation)
	fmt.Print("Enter Alt Pickup Location (optional): ")
	fmt.Scan(&newTrip.AltPickupLocation)
	fmt.Print("Enter Destination: ")
	fmt.Scan(&newTrip.Destination)

	// Consume the newline character from the input buffer
	fmt.Scanln()

	// Prompt for start traveling time as a string using the new function
	newTrip.ScheduledStartTime, _ = getValidDateTimeInput("Enter Start Traveling Time (format: 'yyyy-mm-dd hh:mm'): ")

	fmt.Print("Enter Max Passengers: ")
	fmt.Scan(&newTrip.MaxPassengers)
	newTrip.AvailablePassenger = newTrip.MaxPassengers

	// Insert into the trips table using prepared statements
	stmt, err := db.Prepare("INSERT INTO trips (owner_id, pickup_location, alt_pickup_location, destination, scheduled_start_time, max_passengers, available_passenger) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing SQL statement:", err)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(newTrip.OwnerID, newTrip.PickupLocation, newTrip.AltPickupLocation, newTrip.Destination, newTrip.ScheduledStartTime, newTrip.MaxPassengers, newTrip.AvailablePassenger)
	if err != nil {
		fmt.Println("Error creating trip:", err)
		return
	}

	newTripID, _ := result.LastInsertId()
	newTrip.ID = int(newTripID)

	fmt.Println("Trip created successfully:")
	printTrip(newTrip)
}

// Function to get a valid datetime input from the user
func getValidDateTimeInput(prompt string) (time.Time, error) {
	var input string
	var parsedTime time.Time
	var err error

	possibleLayouts := []string{"2006-01-02 15:04", "2006-01-02 03:04PM"}

	for {
		fmt.Print(prompt)

		// Use bufio.NewReader to handle input more reliably
		reader := bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Try to parse the input string into a time.Time variable using possible layouts
		for _, layout := range possibleLayouts {
			parsedTime, err = time.Parse(layout, input)
			if err == nil {
				break
			}
		}

		if err == nil {
			break // Break out of the loop if parsing is successful
		}

		// If parsing fails, print an error message and ask the user to try again
		fmt.Println("Invalid datetime format. Please use the format 'yyyy-mm-dd HH:mm'.")
	}

	return parsedTime, nil
}

// Function to check for date and time conflicts with enrolled trips
func hasDateConflict(currentUserID int, newTripStartTime time.Time) (bool, error) {
	// Retrieve the scheduled start times of trips joined by the user
	rows, err := db.Query("SELECT t.scheduled_start_time FROM trips t JOIN user_trips ut ON t.id = ut.trip_id WHERE ut.user_id = ?", currentUserID)
	if err != nil {
		return true, err
	}
	defer rows.Close()

	// Check for date and time conflicts with existing trips
	for rows.Next() {
		var enrolledTripStartTimeStr string
		err := rows.Scan(&enrolledTripStartTimeStr)
		if err != nil {
			return true, err
		}

		enrolledTripStartTime, err := time.Parse("2006-01-02 15:04:05", enrolledTripStartTimeStr)
		if err != nil {
			return true, err
		}

		// Check for date and time conflicts
		if newTripStartTime.Before(enrolledTripStartTime.Add(time.Hour)) && newTripStartTime.After(enrolledTripStartTime.Add(-time.Hour)) {
			return true, nil
		}
	}

	return false, nil
}

// Function to get and display available trips for the user to choose from
func getTrip(currentUserID int) {
	// Display all trips
	rows, err := db.Query("SELECT id, owner_id, pickup_location, alt_pickup_location, destination, scheduled_start_time, max_passengers, available_passenger FROM trips")
	if err != nil {
		fmt.Println("Error getting trips:", err)
		return
	}
	defer rows.Close()

	fmt.Println("Available Trips:")
	for rows.Next() {
		var trip Trip

		var scheduledStartTimeStr string // Use a string to temporarily store the scheduled_start_time

		err := rows.Scan(&trip.ID, &trip.OwnerID, &trip.PickupLocation, &trip.AltPickupLocation,
			&trip.Destination, &scheduledStartTimeStr, &trip.MaxPassengers, &trip.AvailablePassenger,
		)
		if err != nil {
			fmt.Println("Error scanning trip:", err)
			return
		}

		// Parse the scheduled_start_time string into a time.Time variable
		scheduledStartTime, err := time.Parse("2006-01-02 15:04:05", scheduledStartTimeStr)
		if err != nil {
			fmt.Println("Error parsing start time:", err)
			return
		}

		// Assign the parsed time to the trip structure
		trip.ScheduledStartTime = scheduledStartTime

		printTrip(trip)
	}

	// Ask the user to choose a trip
	var tripID int
	fmt.Print("Enter Trip ID to book trip (or 0 to go back): ")
	fmt.Scan(&tripID)

	if tripID == 0 {
		return
	}

	// Call the function to get trip details
	getTripDetails(currentUserID, tripID)
}

// Function to get and display details of a selected trip and handle user enrollment
func getTripDetails(currentUserID int, tripID int) {
	if tripID == 0 {
		return
	}

	var trip Trip
	var startTimeStr string

	// Retrieve details of the selected trip
	err := db.QueryRow("SELECT id, owner_id, pickup_location, alt_pickup_location, destination, scheduled_start_time, max_passengers, available_passenger FROM trips WHERE id=?", tripID).Scan(
		&trip.ID, &trip.OwnerID, &trip.PickupLocation, &trip.AltPickupLocation,
		&trip.Destination, &startTimeStr, &trip.MaxPassengers, &trip.AvailablePassenger,
	)

	// Check if the trip exists
	if err == sql.ErrNoRows {
		fmt.Println("Trip not found.")
		return
	} else if err != nil {
		fmt.Println("Error getting trip:", err)
		return
	}

	// Parse the scheduled_start_time string into a time.Time variable
	scheduledStartTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr)
	if err != nil {
		fmt.Println("Error parsing start time:", err)
		return
	}

	// Assign the parsed time to the trip structure
	trip.ScheduledStartTime = scheduledStartTime

	// Check if the user is already part of the trip
	isAlreadyJoined, err := isUserJoinedTrip(currentUserID, tripID)
	if err != nil {
		fmt.Println("Error checking if user is already part of the trip:", err)
		return
	}

	if isAlreadyJoined {
		fmt.Println("You are already part of this trip.")
		return
	}

	// Check if the current user is the owner of the trip
	isOwner, err := isUserTripOwner(currentUserID, tripID)
	if err != nil {
		fmt.Println("Error checking if user is the owner of the trip:", err)
		return
	}

	if isOwner {
		fmt.Println("You cannot book your own trip.")
		return
	}

	// Check if there are available seats in the trip
	if trip.AvailablePassenger > 0 {
		// Check for date and time conflicts
		dateConflict, err := hasDateConflict(currentUserID, trip.ScheduledStartTime)
		if err != nil {
			fmt.Println("Error checking date and time conflicts:", err)
			return
		}

		if dateConflict {
			fmt.Println("You have a date and time conflict with an enrolled trip.")
			return
		}

		// Update the available passengers for the selected trip
		_, err = db.Exec("UPDATE trips SET available_passenger=? WHERE id=?", trip.AvailablePassenger-1, tripID)
		if err != nil {
			fmt.Println("Error updating available passengers:", err)
			return
		}

		// Insert a new record in the user_trips table
		_, err = db.Exec("INSERT INTO user_trips (user_id, trip_id) VALUES (?, ?)", currentUserID, tripID)
		if err != nil {
			fmt.Println("Error inserting record in user_trips table:", err)
			return
		}

		fmt.Println("You have successfully joined the trip!")
	} else {
		fmt.Println("No available seats in this trip.")
	}

	fmt.Println("Trip details:")
	printTrip(trip)
}

// Function to check if the user is already part of a trip
func isUserJoinedTrip(userID int, tripID int) (bool, error) {
	var count int
	// Check if there is a record in the user_trips table with the given user ID and trip ID
	err := db.QueryRow("SELECT COUNT(*) FROM user_trips WHERE user_id=? AND trip_id=?", userID, tripID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Function to update the start time of a trip
func updateStartTime(currentUserID int) {
	var tripID int
	fmt.Print("Enter Trip ID to update start time: ")
	fmt.Scan(&tripID)

	// Check if the logged-in user is the owner of the trip
	isOwner, err := isUserTripOwner(currentUserID, tripID)
	if err != nil {
		fmt.Println("Error checking if user is the owner of the trip:", err)
		return
	}

	if !isOwner {
		fmt.Println("You can only update the start time of your own trips.")
		return
	}

	// Prompt for the updated start time
	var newStartTimeStr string
	fmt.Print("Enter Updated Start Time (format: '2006-01-02 15:04:05', or press Enter for the current time): ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		newStartTimeStr = scanner.Text()
	}

	// Use the current time if the input is empty
	if newStartTimeStr == "" {
		newStartTimeStr = time.Now().Format("2006-01-02 15:04:05")
	}

	// Parse the new start time string into a time.Time variable
	newStartTime, err := time.Parse("2006-01-02 15:04:05", newStartTimeStr)
	if err != nil {
		fmt.Println("Error parsing updated start time:", err)
		return
	}

	// Update the start time of the trip
	_, err = db.Exec("UPDATE trips SET scheduled_start_time=? WHERE id=?", newStartTime, tripID)
	if err != nil {
		fmt.Println("Error updating start time:", err)
		return
	}

	fmt.Println("Start time updated successfully.")
}

// Function to handle the deletion of a trip by the trip owner
func deleteTrip(currentUserID int) {
	var tripID int
	fmt.Print("Enter Trip ID to delete: ")

	// Read input for the trip ID
	_, err := fmt.Scan(&tripID)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Check if the logged-in user is the owner of the trip
	var ownerID int
	var scheduledStartTimeStr string
	err = db.QueryRow("SELECT owner_id, scheduled_start_time FROM trips WHERE id=?", tripID).Scan(&ownerID, &scheduledStartTimeStr)
	if err != nil {
		fmt.Println("Error checking trip owner or retrieving scheduled start time:", err)
		return
	}

	// Only allow the owner to delete their own trips
	if ownerID != currentUserID {
		fmt.Println("You can only delete your own trips.")
		return
	}

	// Parse the scheduled start time string into a time.Time variable
	scheduledStartTime, err := time.Parse("2006-01-02 15:04:05", scheduledStartTimeStr)
	if err != nil {
		fmt.Println("Error parsing scheduled start time:", err)
		return
	}

	// Check if it's within the 30-minute window to allow cancellation
	if !isWithin30MinuteWindow(scheduledStartTime) {
		fmt.Println("You can only cancel trips 30 minutes before the scheduled time.")
		return
	}

	// Continue with the deletion by setting AccountDeleted field
	_, err = db.Exec("DELETE FROM trips WHERE id=?", tripID)

	if err != nil {
		fmt.Println("Error deleting trip:", err)
		return
	}

	fmt.Println("Trip deleted successfully.")
}

// Function to check if it's within the 30-minute window
func isWithin30MinuteWindow(scheduledTime time.Time) bool {
	thirtyMinutesAhead := time.Now().Add(30 * time.Minute)
	return scheduledTime.After(time.Now()) && scheduledTime.Before(thirtyMinutesAhead)
}

// Function to print user details
func printUser(u User) {
	carOwnerStr := "no"
	if u.CarOwner == "yes" {
		carOwnerStr = "yes"
	}
	fmt.Printf("ID: %d, Name: %s %s, Mobile: %s, Email: %s, Car Owner: %s\n", u.ID, u.FirstName, u.LastName, u.MobileNumber, u.EmailAddress, carOwnerStr)
}

// Function to print trip details
func printTrip(t Trip) {
	fmt.Printf("ID: %d, Owner ID: %d, Pickup: %s, Alt Pickup: %s, Destination: %s, Start Time: %s, Max Passengers: %d, Available Passengers: %d\n",
		t.ID, t.OwnerID, t.PickupLocation, t.AltPickupLocation, t.Destination, t.ScheduledStartTime.Format("2006-01-02 15:04:05"), t.MaxPassengers, t.AvailablePassenger)
}
