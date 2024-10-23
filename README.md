# 🚗 Goods Transport booking app

A comprehensive goods transport solution built with modern technologies, featuring real-time tracking, interactive maps, and designed to handle loads of *10k Requests/sec* with **50m users** and **100k drivers**

## 📑 Table of Contents
- [Overview](#overview)
- [Repository Structure](#repository-structure)
- [Key Features](#key-features)
- [Technology Stack](#technology-stack)
- [Getting Started](#getting-started)
- [Branch Information](#branch-information)
- [Contributing](#contributing)
- [License](#license)

## 🌟 Overview

This project implements a complete ride-sharing platform with both frontend and backend components. The platform facilitates real-time ride booking, driver matching, location tracking, and payment processing through a microservices architecture and a modern web interface.

## 📐 Architecture

![architecture](https://github.com/user-attachments/assets/e33a8f4c-597b-4299-ab70-04501e1bb33d)

**You can find the reasoning and explainination behind it [here](https://docs.google.com/document/d/1k2rRMaz8bm45_f6i-2wglWUaDMv8CLnWnFlcAX-j-CE/edit?usp=sharing)**

> ✅ Watch the demo video [here](https://www.youtube.com/watch?v=yAkh492Kj6k)

> ⚠️ You can also view complete system architecture diagram [here](https://excalidraw.com/#json=AAdmRY_7GbA_Niqt3iuET,HN-Y0n9DxFeS7x2KpTAzIg) in exalidraw.

## ✨ Key Features

- Dynamic fare calculation with surge pricing
- Seamless payment processing
- Driver-passenger matching algorithm with locking mechanism
- Interactive map integration
- Real-time location updates for tracking driver
- Comprehensive booking management
- Responsive user interface

## 💻 Technology Stack

### Frontend
- Next.js
- React
- TypeScript
- Tailwind CSS
- shadcn/ui components

### Backend
- Go (Golang)
- PostgreSQL
- Redis
- MongoDB
- RabbitMQ
- Apache Kafka

## 🚀 Getting Started

### Branch Information

This repository maintains three main branches:
- `master` - Main development branch (identical to `web`)
- `web` - Frontend application
- `backend` - Backend services

### Quick Start

1. Clone the repository:
```bash
git clone https://github.com/my-neme-eh-jeff/highly-scalable-goods-transport-app.git app
cd app
```

2. Choose your branch based on your needs:
```bash
# For frontend development
git checkout web

# For backend development
git checkout backend
```

## 📘 Documentation

### Frontend Application (web branch)
The frontend is a Next.js application providing the user interface for the ride-sharing platform. For detailed setup and running instructions, please switch to the `web` branch and refer to its README.

Key frontend features:
- Interactive map for ride booking
- Real-time ride tracking
- Fare estimates and surge pricing
- Responsive design
- Modern UI components

### Backend Services (backend branch)
The backend consists of multiple microservices built with Go. For detailed setup and running instructions, please switch to the `backend` branch and refer to its README.

