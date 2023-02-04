package security

import (
	"crypto/tls"
	"errors"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func ValidateJWT(c *fiber.Ctx) error {
	cert, err := tls.LoadX509KeyPair("./certs/cert.crt", "./keys/key.key")
	if err != nil {
		log.Printf("Failed to load certificates: %s\n", err)
		return errors.New("Request failed")
	}
	tlsConfig := &tls.Config{
		ServerName:   "www.moviesunday.com",
		Certificates: []tls.Certificate{cert},
		// TODO don't skip...
		InsecureSkipVerify: true,
	}
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("https://localhost:4321/validate")
	req.Header.Add("Authorization", c.Get(fiber.HeaderAuthorization))
	req.Header.Add("Username", c.Get("Username"))
	req.Header.SetMethodBytes([]byte("GET"))
	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{
		TLSConfig: tlsConfig,
	}
	err = client.Do(req, resp)
	if err != nil {
		log.Printf("Request failed: %s\n", err)
		return errors.New("Request failed")
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		log.Printf("%s not authorized\n", c.Get("Username"))
		return errors.New(strconv.Itoa(resp.StatusCode()))
	}
	return nil
}
