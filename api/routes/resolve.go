package routes

import (
	"github.com/damilola99-web/url-shortner/database"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")

	redisDb := database.CreateClient(0)
	defer redisDb.Close()
	value, err := redisDb.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "The link not found in the database"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "An error occured on the server please try gain later."})
	}

	rateLimitDb := database.CreateClient(1)
	defer rateLimitDb.Close()
	_ = rateLimitDb.Incr(database.Ctx, "counter")
	return c.Redirect(value, 301)
}
