package routes

import (
	"gihtub.com/ahdirmai/shorten-url-go/database" // Mengimpor paket database.
	"github.com/go-redis/redis/v8"                // Mengimpor paket go-redis.
	"github.com/gofiber/fiber/v2"                 // Mengimpor paket gofiber.
)

// ResolveURL adalah handler untuk mengatasi permintaan pengalihan URL.
func ResolveURL(c *fiber.Ctx) error {
	// Mendapatkan nilai parameter 'url' dari permintaan HTTP.
	url := c.Params("url")
	// Membuat klien Redis dengan nomor database 0 (default).
	r := database.CreateClient(0)
	defer r.Close() // Menutup koneksi Redis setelah selesai menggunakan.

	// Mengambil nilai dari Redis berdasarkan 'url' yang diberikan.
	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		// Jika 'url' tidak ditemukan dalam database Redis, kembalikan respons dengan status 404 (Not Found).
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"errors": "short not found in database",
		})
	} else if err != nil {
		// Jika terjadi kesalahan dalam menghubungkan ke database Redis, kembalikan respons dengan status 500 (Internal Server Error).
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": "Cannot Connect To DB",
		})
	}

	// Membuat klien Redis baru dengan nomor database 1.
	rInr := database.CreateClient(1)
	defer rInr.Close() // Menutup koneksi Redis yang baru dibuat setelah selesai menggunakan.

	// Meningkatkan hitungan di database Redis dengan nomor database 1.
	_ = rInr.Incr(database.Ctx, "counter")

	// Mengarahkan pengguna ke URL asli dengan status pengalihan 301 (Moved Permanently).
	return c.Redirect(value, 301)
}
