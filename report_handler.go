package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// apis.GET("/reports/projects/:projectid", reportOnProject(db))
func reportOnProject(db *Database) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return echo.NewHTTPError(http.StatusInternalServerError, "Not Yet Implemented")
	}
}

// apis.GET("/reports/date/:datestr", reportOnTime(db))
func reportOnTime(db *Database) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return echo.NewHTTPError(http.StatusInternalServerError, "Not Yet Implemented")
	}
}
