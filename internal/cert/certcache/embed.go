// Package certcache provides an embedded cache of certificates.
package certcache

import "embed"

//go:embed root/EC_PK.bin
var root []byte

//go:embed g1/*.bin
var g1 embed.FS

//go:embed g2/*.bin
var g2 embed.FS

// Root returns the ERCA root certificate.
func Root() []byte {
	return root
}

// ReadG1 reads a cached Gen1 certificate by its CHR.
func ReadG1(chr string) ([]byte, bool) {
	data, err := g1.ReadFile(chr + ".bin")
	if err != nil {
		return nil, false
	}
	return data, true
}

// ReadG2 reads a cached Gen2 certificate by its CHR.
func ReadG2(chr string) ([]byte, bool) {
	data, err := g2.ReadFile(chr + ".bin")
	if err != nil {
		return nil, false
	}
	return data, true
}
