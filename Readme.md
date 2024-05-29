# Digital Inventory Management System

## Table of Contents
- [Introduction](#introduction)
- [Technologies Used](#technologies-used)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [API Documentation](#api-documentation)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)

## Introduction
The Digital Inventory Management System is designed to streamline the management of digital assets such as CPUs, RAM, PCs, monitors, VGA cards, and more. This system leverages the power of Clean Architecture and JWT to ensure scalability, security, and efficient management of inventory.

## Technologies Used
- **Golang**: Programming language for backend development.
- **Gin Gonic**: Web framework used to build the RESTful API.
- **GORM**: ORM library for Golang, used for interacting with PostgreSQL.
- **JWT (JSON Web Token)**: Used for secure authentication and authorization.
- **bcrypt**: Password hashing algorithm.
- **PostgreSQL**: Relational database management system.
- **OpenAPI**: For API documentation.

## Features
- Secure user authentication and authorization using JWT.
- Role-based access control for admins and users.
- CRUD operations for managing digital inventory items.
- Detailed inventory reports.
- Activity tracking for auditing purposes.

## Installation

### Prerequisites
- [Go](https://golang.org/doc/install) (version 1.18+)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Git](https://git-scm.com/)

### Steps
1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/digital-inventory.git
   cd digital-inventory
   
2. **Set up PostgreSQL database:**
   ```bash
   CREATE DATABASE inventory_db;

3. **Install dependencies:**
   ```bash
   go mod tidy

4. **Run the application:**
   ```bash
   go run main.go
   
### Usage
