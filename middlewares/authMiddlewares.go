package middlewares

import (
	"bougette-backend/common"
	"bougette-backend/helper"
	"strings"

	"github.com/labstack/echo/v4"
)

func IsAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// By using the authentication in Vary, the server ensures that the caching mechanism does
		// not cache responses for different authentication states.
		// Meant for authorized users and serve them to unauthenticated users or vice versa.
		// It's essentially prevents security risks related to caching.
		c.Response().Header().Add("Vary", "Authorization")
		authHeader := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return common.SendUnauthorizedResponse(c, "Unauthorized: Missing or invalid token")
		}

		accessToken := strings.TrimPrefix(authHeader, "Bearer ")
		if accessToken == "" {
			return common.SendUnauthorizedResponse(c, "Authorization token is empty")
		}

		claims, err := helper.ParseAccessToken(accessToken)
		if err != nil {
			return common.SendUnauthorizedResponse(c, "Unauthorized")
		}

		if helper.IsTokenExpired(claims) {
			return common.SendUnauthorizedResponse(c, "Unauthorized: Token expired")
		}

		c.Set("userID", claims.UserID)
		return next(c)
	}
}
