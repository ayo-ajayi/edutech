# Edutech API Documentation

This readme provides an overview of the Edutech API, detailing the project structure, dependencies, and key functionalities.

## Table of Contents
- [Edutech API Documentation](#edutech-api-documentation)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Setup and Configuration](#setup-and-configuration)
  - [API Endpoints](#api-endpoints)
  - [Authentication and Authorization](#authentication-and-authorization)
  - [Middleware](#middleware)
  - [Health Check](#health-check)
  - [License](#license)
  - [Author](#author)



## Introduction

The Edutech API is designed to manage students, tutors, subjects, and authentication within an educational platform. It utilizes the Gin framework for building the API, MongoDB for data storage, and various external services.


## Setup and Configuration

1. Clone the repository:
   ```bash
   git clone https://github.com/ayo-ajayi/edutech.git
    ```
2. Install dependencies:
   ```bash
   go get .
   ```
3. Create a `.env` file in the root directory of the project and add the following environment variables:
- `MONGODB_URI`: MongoDB connection URI
- `MONGODB_NAME`: MongoDB database name
- `EMAIL_API_KEY`: API key for sending emails
- `EMAIL_SENDER_NAME`: Sender name for outgoing emails
- `EMAIL_SENDER_ADDRESS`: Sender email address
- `BASE_URL`: Base URL for email verification links
- `ACCESS_TOKEN_SECRET`: Secret key for JWT token generation

4. Run the application:
   ```bash
   go run main.go
   ```
   
5. The application will be available at `http://localhost:8000`

## API Endpoints

- **POST** `/api/v1/login`: User login
- **POST** `/api/v1/forgot-password`: Request to reset password
- **POST** `/api/v1/reset-password`: Reset user password
- **GET** `/api/v1/verify/:token`: Verify user email
- **DELETE** `/api/v1/logout`: User logout
- **POST** `/api/v1/students`: Student registration
- **GET** `/api/v1/students/profile`: Get student profile
- **GET** `/api/v1/students/subjects`: Get registered subjects for a student
- **POST** `/api/v1/students/subjects`: Register a subject for a student
- **POST** `/api/v1/students/tutors/register`: Register a tutor for a student
- **GET** `/api/v1/students/tutors`: Get registered tutors for a student
- **POST** `/api/v1/tutors`: Tutor registration
- **GET** `/api/v1/tutors/profile`: Get tutor profile
- **POST** `/api/v1/subjects`: Create a new subject

## Authentication and Authorization

- **Authentication:** JWT (JSON Web Tokens) is used for user authentication.
- **Authorization:** Middleware ensures that only authenticated users with the correct role can access specific endpoints.

## Middleware

The application uses middleware for common functionality:

- `jsonMiddleware`: Sets the content type to JSON for all responses.
- `auth.NewCors()`: Adds CORS headers to allow cross-origin requests.

## Health Check

The `/healthz` endpoint provides a simple health check response, indicating the API is operational.


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author
- [Ayomide Ajayi](https://github.com/ayo-ajayi)