package kiviasdk

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
)

// FiberLog returns Fiber v3 middleware that asynchronously logs request metrics.
// Register it before your routes with app.Use(client.FiberLog()).
// Downstream errors are handled by the app's error handler before capturing status.
func (c *Client) FiberLog() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		start := time.Now()
		err := ctx.Next()
		if err != nil {
			err = ctx.App().ErrorHandler(ctx, err)
		}
		// Fiber reuses request memory; copy strings before starting the goroutine.
		entry := &Log{
			Path:      strings.Clone(ctx.OriginalURL()),
			Status:    ctx.Response().StatusCode(),
			IPAddress: strings.Clone(ctx.IP()),
			Timestamp: time.Now(),
			Latency:   int(time.Since(start).Milliseconds()),
		}
		go c.sendLogEntry(entry)
		return err
	}
}
