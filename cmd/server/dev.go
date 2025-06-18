//go:build dev

package main

import "os"

func init() {
	os.Setenv("HTTP_PORT", "8080")
}
