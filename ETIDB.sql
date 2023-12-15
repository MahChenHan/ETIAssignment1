-- Create the database
CREATE DATABASE IF NOT EXISTS my_db;

-- Use the database
USE my_db;

-- Create the Users table
CREATE TABLE IF NOT EXISTS Users (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    mobile_number VARCHAR(20),
    email_address VARCHAR(255),
    driver_license VARCHAR(255),
    car_plate_number VARCHAR(20),
    car_owner VARCHAR(3),
    account_created TIMESTAMP,
    last_updated TIMESTAMP,
    account_deleted TIMESTAMP
);
-- Create the Trips table
CREATE TABLE IF NOT EXISTS trips (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    owner_id INT,
    pickup_location VARCHAR(255),
    alt_pickup_location VARCHAR(255),
    destination VARCHAR(255),
    max_passengers INT,
    available_passenger INT,
	scheduled_start_time DATETIME
);

CREATE TABLE IF NOT EXISTS user_trips (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    trip_id INT,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (trip_id) REFERENCES trips(id)
);

Select * from user_trips;
select * from Trips;
select * from users;