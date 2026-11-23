package config

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
)

type Config struct {
	Puerto, Ambiente, Base, Frontend, JWKS, Emisor, Audiencia, Cookie, ClaveRuta string
	Mocks                                                                        bool
}

func Leer() (Config, error) {
	c := Config{Puerto: os.Getenv("API_PORT"), Ambiente: os.Getenv("API_ENV"), Base: os.Getenv("DATABASE_URL"), Frontend: os.Getenv("FRONTEND_URL"), JWKS: os.Getenv("AUTH_JWKS_URL"), Emisor: os.Getenv("AUTH_ISSUER"), Audiencia: os.Getenv("AUTH_AUDIENCE"), Cookie: os.Getenv("AUTH_COOKIE_NAME"), ClaveRuta: os.Getenv("SIGNING_KEY_FILE"), Mocks: os.Getenv("API_MOCKS") == "true"}
	for k, v := range map[string]string{"API_PORT": c.Puerto, "API_ENV": c.Ambiente, "DATABASE_URL": c.Base, "FRONTEND_URL": c.Frontend, "AUTH_JWKS_URL": c.JWKS, "AUTH_ISSUER": c.Emisor, "AUTH_AUDIENCE": c.Audiencia, "AUTH_COOKIE_NAME": c.Cookie, "SIGNING_KEY_FILE": c.ClaveRuta, "TZ": os.Getenv("TZ")} {
		if v == "" {
			return c, fmt.Errorf("falta configuración %s", k)
		}
	}
	if c.Puerto != "8080" {
		return c, fmt.Errorf("API_PORT debe ser 8080")
	}
	if c.Ambiente != "dev" && c.Ambiente != "test" && c.Ambiente != "prod" {
		return c, fmt.Errorf("API_ENV inválido")
	}
	if c.Mocks && c.Ambiente == "prod" {
		return c, fmt.Errorf("mocks prohibidos en prod")
	}
	for _, v := range []string{c.JWKS, c.Frontend} {
		u, e := url.Parse(v)
		if e != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
			return c, fmt.Errorf("URL inválida")
		}
		if c.Ambiente == "prod" && u.Scheme != "https" {
			return c, fmt.Errorf("HTTPS obligatorio en prod")
		}
	}
	return c, nil
}
func LeerClave(ruta string) (ed25519.PrivateKey, error) {
	b, err := os.ReadFile(ruta)
	if err != nil {
		return nil, fmt.Errorf("leer clave de firma: %w", err)
	}
	clave, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil || len(clave) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("clave Ed25519 inválida")
	}
	esperada := ed25519.NewKeyFromSeed(clave[:ed25519.SeedSize])
	if !ed25519.PrivateKey(clave).Equal(esperada) {
		return nil, fmt.Errorf("clave Ed25519 inconsistente")
	}
	return ed25519.PrivateKey(clave), nil
}
