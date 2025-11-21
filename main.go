package main

import (
	"fmt"
	"os"
	"student-performance-report/config"
	
)

func main() {
    config.LoadEnv() // Load file .env

    host := os.Getenv("DB_HOST")

    if host == "" {
        fmt.Println(".env gagal diload atau DB_HOST tidak ditemukan")
    } else {
        fmt.Println(".env berhasil diload. DB_HOST =", host)
    }
}
