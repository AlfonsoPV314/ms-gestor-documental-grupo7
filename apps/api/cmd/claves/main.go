package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	ruta := os.Getenv("SIGNING_KEY_FILE")
	if ruta == "" {
		ruta = ".secrets/firma.key"
	}
	if _, err := os.Stat(ruta); err == nil {
		fmt.Println("Clave existente conservada")
		return
	}
	if err := os.MkdirAll(filepath.Dir(ruta), 0700); err != nil {
		panic(err)
	}
	_, privada, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(ruta, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		panic(err)
	}
	if _, err = f.WriteString(base64.StdEncoding.EncodeToString(privada)); err != nil {
		panic(err)
	}
	if err = f.Close(); err != nil {
		panic(err)
	}
	fmt.Println("Clave de firma creada; conservarla junto con los respaldos")
}
