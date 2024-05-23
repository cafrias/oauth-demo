package auth

import (
	"auth-server/internal/common"
	"fmt"

	"github.com/labstack/echo/v4"
)

// Authenticated is a middleware that checks if the user is authenticated
func Authenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.(common.AppContext)
		s, _ := ctx.GetSession()

		if !s.IsAuthenticated() {
			url := fmt.Sprintf("/login?redirect=%s", c.Request().URL)
			return c.Redirect(302, url)
		}

		return next(ctx)
	}
}
