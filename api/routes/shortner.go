package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/damilola99-web/url-shortner/database"
	"github.com/damilola99-web/url-shortner/helpers"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"customLink"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL            string        `json:"url"`
	CustomShort    string        `json:"customLink"`
	Expiry         time.Duration `json:"expiry"`
	XRateRemaining int           `json:"xRateRemaining"`
}

func ShortenURL(c *fiber.Ctx) error {
	body := new(request)
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse json"})
	}

	// rate limiting

	rateLimitDb := database.CreateClient(1)
	defer rateLimitDb.Close()
	value, err := rateLimitDb.Get(database.Ctx, c.IP()).Result()
	_ = value
	if err == redis.Nil {
		_ = rateLimitDb.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()

	} else {
		value, _ := rateLimitDb.Get(database.Ctx, c.IP()).Result()
		intValue, _ := strconv.Atoi(value)
		if intValue <= 0 {
			limit, _ := rateLimitDb.TTL(database.Ctx, c.IP()).Result()
			_ = limit
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Rate limit exceeded."})
		}
	}

	// check if input is a url
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL sent."})
	}

	// check for domain error
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Stop playing I already caught you!! You wan hack my API?"})
	}

	// enforce ssl
	body.URL = helpers.EnforceHTTP(body.URL)

	var id string

	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	linkDb := database.CreateClient(0)
	defer linkDb.Close()

	value, _ = linkDb.Get(database.Ctx, id).Result()
	if value != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "URL custom short is already used."})
	}
	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = linkDb.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server error."})
	}

	resp := response{
		URL:            body.URL,
		CustomShort:    "",
		Expiry:         body.Expiry,
		XRateRemaining: 10,
	}

	rateLimitDb.Decr(database.Ctx, c.IP())

	value, _ = rateLimitDb.Get(database.Ctx, c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(value)
	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id
	return c.Status(fiber.StatusOK).JSON(resp)
}
