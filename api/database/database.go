package database

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

// Membuat konteks (context) yang akan digunakan dalam komunikasi dengan Redis.
var Ctx = context.Background()

// CreateClient adalah fungsi yang digunakan untuk membuat klien Redis.
// Fungsi ini menerima parameter dbNo yang merupakan nomor database yang akan digunakan.
func CreateClient(dbNo int) *redis.Client {
	// Membuat klien Redis dengan menggunakan konfigurasi yang diberikan.
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDR"), // Mengambil alamat server Redis dari variabel lingkungan DB_ADDR.
		Password: os.Getenv("DB_PASS"), // Mengambil kata sandi Redis dari variabel lingkungan DB_PASS.
		DB:       dbNo,                 // Menggunakan nomor database yang diberikan sebagai dbNo.
	})

	// Mengembalikan klien Redis yang telah dibuat.
	return rdb
}
