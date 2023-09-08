package helpers

import (
	"os"
	"strings"
)

// EnforceHTTP adalah fungsi yang memastikan URL dimulai dengan "http://" atau "https://".
func EnforceHTTP(url string) string {
	// Memeriksa apakah URL tidak dimulai dengan "http" atau "https".
	if url[:4] != "http" {
		// Jika tidak, tambahkan "http://" pada awal URL.
		return "http://" + url
	}
	// Jika sudah dimulai dengan "http" atau "https", kembalikan URL asli.
	return url
}

// RemoveDomainError adalah fungsi yang memeriksa apakah URL sama dengan nilai variabel lingkungan "DOMAIN".
func RemoveDomainError(url string) bool {
	// Memeriksa apakah URL sama dengan nilai "DOMAIN" yang tersimpan dalam variabel lingkungan.
	if url == os.Getenv("DOMAIN") {
		// Jika sama, kembalikan false untuk menandakan ada kesalahan.
		return false
	}

	// Menghapus awalan "http://" atau "https://", dan "www." dari URL.
	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)

	// Memisahkan URL berdasarkan tanda '/' dan mengambil bagian pertama (domain).
	newURL = strings.Split(newURL, "/")[0]

	// Memeriksa kembali apakah URL yang sudah dimodifikasi sama dengan "DOMAIN".
	if newURL == os.Getenv("DOMAIN") {
		// Jika sama, kembalikan false untuk menandakan ada kesalahan.
		return false
	}

	// Jika URL tidak sama dengan "DOMAIN" setelah modifikasi, kembalikan true.
	return true
}
