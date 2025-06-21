//go:build dev

package main

import "os"

func init() {
	os.Setenv("HTTP_PORT", "8080")
	os.Setenv("API_KEY", "awesome-api-key")
}
