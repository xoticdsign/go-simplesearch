package utils

import "fmt"

// BuildAddress() combines a host and port into a single address string.
//
// This function takes two input parameters: host and port, and returns a formatted
// string in the form of 'host:port'. It is useful for constructing an address
// that can be used to configure network connections (e.g., HTTP server or client).
func BuildAddress(host string, port string) string {
	return fmt.Sprintf("%s:%s", host, port)
}
