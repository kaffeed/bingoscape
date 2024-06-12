package handlers

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/kaffeed/bingoscape/views/errors"
	"github.com/labstack/echo/v4"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)

	var errorPage func(fp bool) templ.Component

	switch code {
	case 401:
		errorPage = errors.Error401
	case 404:
		errorPage = errors.Error404
	case 500:
		errorPage = errors.Error500
	}

	isAuthenticated, _ := c.Get("ISAUTHENTICATED").(bool)
	isManagement, _ := c.Get(mgmnt_key).(bool)
	// isError = true
	c.Set("ISERROR", true)

	render(c, errors.ErrorIndex(
		fmt.Sprintf("| Error (%d)", code),
		"",
		isAuthenticated,
		isManagement,
		c.Get("ISERROR").(bool),
		errorPage(isAuthenticated),
	))
}
