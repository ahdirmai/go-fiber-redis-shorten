package routes

import (
	// Mengimpor berbagai paket yang diperlukan.
	"os"
	"strconv"
	"time"

	"gihtub.com/ahdirmai/shorten-url-go/database"
	"gihtub.com/ahdirmai/shorten-url-go/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Definisikan struktur 'request' untuk menangani payload JSON yang masuk.
type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

// Definisikan struktur 'response' untuk menangani respons JSON yang keluar.
type response struct {
	URL            string        `json:"url"`
	CustomShort    string        `json:"short"`
	Expiry         time.Duration `json:"expiry"`
	XRateRemaining int           `json:"rate_limit"`
	XRateLimitRest time.Duration `json:"rate_limit_reset"`
}

// ShortenURL adalah handler untuk memendekkan URL.
func ShortenURL(c *fiber.Ctx) error {
	// Parse payload JSON yang masuk ke dalam struktur 'request'.
	body := new(request)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Implementasi Rate Limit
	r2 := database.CreateClient(1)
	defer r2.Close()

	// Memeriksa apakah pengguna telah melebihi batasan rate limit.
	val, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		val, _ := r2.Get(database.Ctx, c.IP()).Result()
		valInt, _ := strconv.Atoi(val)

		if valInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":           "Rate limit exceeded",
				"rate_limit_rest": limit / time.Nanosecond / time.Minute,
			})
		}
	}

	// Memeriksa apakah input adalah URL yang valid.
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL",
		})
	}

	// Memeriksa apakah terjadi kesalahan domain pada URL.
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "You can't hack the system",
		})
	}

	// Memastikan URL dimulai dengan "http://" atau "https://".
	body.URL = helpers.EnforceHTTP(body.URL)

	// Membuat ID pendek untuk URL jika tidak ada ID khusus yang diberikan.
	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	// Membuat koneksi Redis ke database 0.
	r := database.CreateClient(0)
	defer r.Close()

	// Memeriksa apakah ID yang telah digunakan sebelumnya.
	val, _ = r.Get(database.Ctx, id).Result()
	if val != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "CustomShort URL is already in use",
		})
	}

	// Mengatur nilai URL ke dalam Redis dengan ID sebagai kunci dan mengatur waktu kadaluwarsa (expiry).
	if body.Expiry == 0 {
		body.Expiry = 24 // Jika expiry tidak ditentukan, maka defaultnya 24 jam.
	}
	err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to connect to server",
		})
	}

	// Mengisi respons yang akan dikirimkan kembali kepada pengguna.
	resp := response{
		URL:            body.URL,
		CustomShort:    "",
		Expiry:         body.Expiry,
		XRateRemaining: 10,
		XRateLimitRest: 30,
	}
	r2.Decr(database.Ctx, c.IP()) // Mengurangi hitungan rate limit untuk pengguna yang sedang memanggil.

	val, _ = r2.Get(database.Ctx, c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	tll, _ := r2.TTL(database.Ctx, c.IP()).Result()
	resp.XRateLimitRest = tll / time.Nanosecond / time.Minute
	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id // Membuat URL pendek dengan domain.

	// Mengirimkan respons JSON ke pengguna dengan status OK.
	return c.Status(fiber.StatusOK).JSON(resp)
}
