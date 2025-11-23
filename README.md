Student Performance Reporting System (Backend)

📋 Overview

This is a backend REST API service designed to manage and report student achievements. The system implements a Dual-Database Architecture (Hybrid SQL & NoSQL) to leverage the relational integrity of PostgreSQL for user management and workflows, while utilizing the flexibility of MongoDB for storing dynamic achievement details (competitions, seminars, organizations, etc.).

The application follows Clean Architecture principles and implements strictly typed Role-Based Access Control (RBAC) for Students, Lecturers (Advisors), and Admins.

🚀 Key Features

Dual-Write Mechanism: Synchronized data storage between PostgreSQL (References & Status) and MongoDB (Dynamic Details).

RBAC Middleware: Secure endpoints with role validation (Mahasiswa, Dosen Wali, Admin) and permission checks.

Achievement Workflow: Full lifecycle management: Draft -> Submitted -> Verified / Rejected.

Soft Delete: Safe deletion strategy for data integrity.

File Attachments: Support for uploading evidence/certificates.

Reporting & Analytics: Aggregated statistics using MongoDB Pipelines and PostgreSQL joins.

Clean Code Structure: Separation of concerns (Handler -> Service -> Repository).

🛠️ Tech Stack

Language: Go (Golang)

Framework: Fiber v2

Databases:

PostgreSQL: Used for Users, Roles, Students, Lecturers, and Achievement Status Tracking.

MongoDB: Used for storing complex and varying achievement detail schemas.

Authentication: JWT (JSON Web Token)

Drivers: lib/pq (Raw SQL), mongo-driver.

📂 Project Structure

student-performance-report
├── app       
│   ├── models           # Data Structures / Entities
│   │   ├── mongodb      # Structs for MongoDB collections (Achievement details)
│   │   └── postgresql   # Structs for SQL tables (Users, Roles, References)
│   ├── repository       # Data Access Layer (Database Queries)
│   │   ├── mongodb      # Implementation of MongoDB operations
│   │   └── postgresql   # Implementation of PostgreSQL operations
│   └── service          # Business Logic Layer
│       ├── mongodb      # Services handling MongoDB logic (Achievement, Reports)
│       └── postgresql   # Services handling SQL logic (Auth, Admin, Student)     
├── config               # Configuration setup (Env, JWT)
├── database             # Connection logic for MongoDB & PostgreSQL
├── docs                 # API Documentation files
├── fiber                # Fiber app configuration
├── middleware           # Auth & Role-Based Access Control (RBAC)
├── pwhash               # Password hashing utilities
├── route                # API Endpoint definitions
├── uploads              # Directory for static file storage (attachments)
├── utils                # Helper functions (Token generators, Validators)
├── .env                 # Environment variables configuration
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
└── main.go              # Application entry point


⚙️ Installation & Setup

Prerequisites

Go 1.20 or higher

PostgreSQL

MongoDB

1. Clone the Repository

git clone [https://github.com/YOUR_USERNAME/student-performance-report.git](https://github.com/YOUR_USERNAME/student-performance-report.git)
cd student-performance-report


2. Install Dependencies

go mod tidy


3. Environment Variables

Create a .env file in the root directory:

PORT=8080

# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=yourpassword
DB_NAME=student_performance_db

# MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DB_NAME=student_performance_mongo

# JWT
JWT_SECRET=your_super_secret_key


4. Database Migration

Execute the provided SQL scripts (in database/migrations if available) to create tables: users, roles, students, lecturers, achievement_references, etc.

5. Run the Application

go run main.go


🔌 API Endpoints

5.1 Authentication

Method

Endpoint

Description

POST

/api/v1/auth/login

User login & JWT generation

POST

/api/v1/auth/refresh

Refresh access token

POST

/api/v1/auth/logout

Logout (Invalidate session)

GET

/api/v1/auth/profile

Get current user profile

5.2 Users (Admin)

Method

Endpoint

Description

Access

GET

/api/v1/users

List all users

Admin

GET

/api/v1/users/:id

Get user detail by ID

Admin

POST

/api/v1/users

Create new user

Admin

PUT

/api/v1/users/:id

Update user data

Admin

DELETE

/api/v1/users/:id

Delete user

Admin

PUT

/api/v1/users/:id/role

Assign role to user

Admin

5.4 Achievements

Method

Endpoint

Description

Access

GET

/api/v1/achievements

List achievements (Filtered by role: own, advisees, or all)

All

GET

/api/v1/achievements/:id

Get achievement detail

All

POST

/api/v1/achievements

Create new achievement (Draft)

Mahasiswa

PUT

/api/v1/achievements/:id

Update achievement (Draft only)

Mahasiswa

DELETE

/api/v1/achievements/:id

Soft delete achievement (Draft only)

Mahasiswa

POST

/api/v1/achievements/:id/submit

Submit achievement for verification

Mahasiswa

POST

/api/v1/achievements/:id/verify

Verify/Approve achievement

Dosen Wali

POST

/api/v1/achievements/:id/reject

Reject achievement with note

Dosen Wali

GET

/api/v1/achievements/:id/history

View achievement status history

All

POST

/api/v1/achievements/:id/attachments

Upload supporting files/certificates

Mahasiswa

5.5 Students & Lecturers

Method

Endpoint

Description

Access

GET

/api/v1/students

List all students

Authorized

GET

/api/v1/students/:id

Get student profile detail

Authorized

GET

/api/v1/students/:id/achievements

Get achievements of a specific student

Authorized

PUT

/api/v1/students/:id/advisor

Assign academic advisor to student

Authorized

GET

/api/v1/lecturers

List all lecturers

Authorized

GET

/api/v1/lecturers/:id/advisees

Get list of advisees (mahasiswa bimbingan)

Lecturer/Admin

5.8 Reports & Analytics

Method

Endpoint

Description

Access

GET

/api/v1/reports/statistics

Global achievement statistics

Admin

GET

/api/v1/reports/student/:id

Specific student performance report

Admin/Lecturer/Owner

🛡️ Security

JWT Authentication: All protected routes require a valid Bearer Token.

Role-Based Middleware: Ensures Students can only modify their own data, and Lecturers can only verify their advisees.

Input Validation: Strict struct binding and validation.

🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

📄 License

This project is licensed under the MIT License.
