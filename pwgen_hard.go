package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"log"
)

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	// https://www.keylength.com/en/compare/
	const keyLength = 24 // 192 bits of entropy

	for i := 0; i < 10; i++ {
		buf := make([]byte, keyLength)

		n, err := rand.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		if n != keyLength {
			log.Fatalf("unexpected number of bytes read: %s", n)
		}

		log.Printf("password: %s", base64.RawStdEncoding.EncodeToString(buf))
	}
}
