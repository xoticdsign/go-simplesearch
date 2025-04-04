package utils

import "fmt"

// Build an address from two parts: host+port.
func BuildAddress(host string, port string) string {
	return fmt.Sprintf("%s:%s", host, port)
}
