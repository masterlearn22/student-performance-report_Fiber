package main

import (
	// "fmt"
	// "os"
	// "context"
	"student-performance-report/config"
	"student-performance-report/database"
	FiberApp "student-performance-report/fiber"

	
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
	defer database.DB.Close()

	// Connect to MongoDB
	database.ConnectMongo()

	//3 Setup Fiber App
	FiberApp.SetupFiber()
	

	//4. Setup Route
}
