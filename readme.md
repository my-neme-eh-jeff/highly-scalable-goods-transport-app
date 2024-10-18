# Ride-Sharing Platform Backend Services 🚗

A robust ride-sharing platform built with Go microservices architecture, designed for scalability and real-time operations.

## 📑 Table of Contents
- [Introduction](#introduction)
- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Setting Up the Environment](#setting-up-the-environment)
- [Running the Backend Services](#running-the-backend-services)
- [Services Overview](#services-overview)
- [Additional Notes](#additional-notes)

## 🚀 Introduction

This project implements a comprehensive ride-sharing platform consisting of multiple backend microservices built with Go (Golang). The platform includes:

- Payment Service
- Transport Service
- Transport Tracking Service
- Driver Location Service
- Transport Matcher Service

## 📋 Prerequisites

Ensure you have the following installed on your system:

- **Go (Golang)** 1.16 or later - [Download Go](https://golang.org/dl/)

### Required cloud hosted urls for these databases and message brokers:

- **PostgreSQL** - [Download PostgreSQL](https://www.postgresql.org/download/)
- **Redis** - [Download Redis](https://redis.io/download)
- **MongoDB** - [Download MongoDB](https://www.mongodb.com/try/download/community)
- **RabbitMQ** - [Download RabbitMQ](https://www.rabbitmq.com/download.html)
- **Apache Kafka** - [Download Kafka](https://kafka.apache.org/downloads)

> **Note**: You can alternatively use setlf hosting for these databases and message brokers. Make sure you have the connection strings and credentials ready.

## 📁 Project Structure

```
backend/
├── payment_service/
├── transport_service/
├── transport_tracking_service/
├── driver_location_service/
├── transport_matcher_service/
├── start_services.sh
└── README.md
```

## 🛠 Setting Up the Environment

### 1. Clone the Repository

```bash
git clone https://github.com/your_username/your_repository.git
cd your_repository
```

### 2. Set Environment Variables

Create a `.env.local` file in each directory according to the .env.example given. You will require:

```bash
# Redis connection string
REDIS_UPSTASH_ADDR="your_redis_connection_string"

# PostgreSQL connection string
NEON_DB_URL="your_postgres_connection_string"

# MongoDB connection string
MONGODB_URI="your_mongodb_connection_string"

# RabbitMQ connection string
CLOUDAMQP_URL="your_rabbitmq_connection_string"

# Kafka connection properties will go in client.properties in root of "transport_service" and "transport_matcher_services"
```

> ⚠️ **Important**: Do not commit the `.env.local` file to version control!

### 3. Install Go Dependencies

Navigate to each service directory and run:

```bash
go mod init
go mod tidy
```

## 🚀 Running the Backend Services

### Using the start_services.sh Script

1. Make the script executable:
```bash
chmod +x start_services.sh
```

2. Run the script:
```bash
./start_services.sh
```

The script will:
- Navigate to each service directory
- Initialize Go modules
- Install dependencies
- Run services in background
- Create and manage log files

### Manual Service Startup

For each service, follow these steps:

```bash
cd backend/service_name
go mod init service_name
go mod tidy
go run .
```

## 🔧 Services Overview

### 1. Payment Service
- **Directory**: `backend/payment_service`
- **Port**: 8080
- **Features**: Fare calculation, surge pricing, payment records

### 2. Transport Service
- **Directory**: `backend/transport_service`
- **Port**: 8081
- **Features**: Transport bookings, driver responses

### 3. Transport Tracking Service
- **Directory**: `backend/transport_tracking_service`
- **Port**: 8082
- **Features**: Real-time transport tracking via SSE

### 4. Driver Location Service
- **Directory**: `backend/driver_location_service`
- **Port**: 8083
- **Features**: Real-time driver location updates via WebSockets

### 5. Transport Matcher Service
- **Directory**: `backend/transport_matcher_service`
- **Port**: 8084
- **Features**: User-driver matching algorithms

## 📝 Additional Notes

### Database Management
- Services auto-create required tables
- Consider using migrations in production
- Regular backup recommendations

### Error Handling
- Comprehensive error handling implemented
- Error reporting via structured logs
- Monitoring endpoints available

### Dependencies
- Managed via Go modules
- Auto-installed via `go mod tidy`
- Version controls in place

### Kafka Configuration
- Requires `client.properties` file
- Separate configs per service
- Secure connection settings
