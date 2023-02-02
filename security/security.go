package security

import (
	"errors"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func ValidateJWT(c *fiber.Ctx) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI("http://localhost:4321/validate")
	req.Header.Add("Authorization", c.Get(fiber.HeaderAuthorization))
	req.Header.Add("Username", c.Get("Username"))
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := fasthttp.Do(req, resp)
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
