package main

import (
	"fmt"
	// "os"
	// "context"
	"student-performance-report/config"
	"student-performance-report/database"
	FiberApp "student-performance-report/fiber"
	routePostgre "student-performance-report/route/postgresql"
	routeMongo "student-performance-report/route/mongodb"

	
)

func main() {

	// 1. Load .env file
    config.LoadEnv() // Load file .env
    // host := os.Getenv("DB_HOST")
    // if host == "" {
    //     fmt.Println(".env gagal diload atau DB_HOST tidak ditemukan")
    // } else {
    //     fmt.Println(".env berhasil diload. DB_HOST =", host)
    // }

	//2. Connect to Database

	// Connect to PostgreSQL
	database.ConnectPostgres()
	defer database.PostgresDB.Close()

	// Connect to MongoDB
	database.ConnectMongo()

	//3 Setup Fiber App
	app := FiberApp.SetupFiber()

	//4. Setup Route
	routePostgre.SetupPostgresRoutes(app, database.PostgresDB)
	routeMongo.SetupMongoRoutes(app, database.MongoDB)

	fmt.Println("Setup route berhasil")
}
