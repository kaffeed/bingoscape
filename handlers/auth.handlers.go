package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/kaffeed/bingoscape/db"
	authviews "github.com/kaffeed/bingoscape/views/auth"
	components "github.com/kaffeed/bingoscape/views/components"
	"golang.org/x/crypto/bcrypt"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const (
	auth_sessions_key string = "authenticate-sessions"
	auth_key          string = "authenticated"
	user_id_key       string = "user_id"
	username_key      string = "username"
	tzone_key         string = "time_zone"
	mgmnt_key         string = "is_management"
)

/********** Handlers for Auth Views **********/

type AuthService interface {
	CreateUser(params db.CreateLoginParams) error
	CheckUsername(username string) (db.Login, error)
	GetAllUsers() ([]db.Login, error)
	DeleteUser(uid int32) error
	UpdatePassword(uid int32, p string) (db.Login, error)
}

func NewAuthHandler(us AuthService) *AuthHandler {
	return &AuthHandler{
		UserServices: us,
	}
}

type AuthHandler struct {
	UserServices AuthService
}

func (ah *AuthHandler) handleUsermanagement(c echo.Context) error {
	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	if !ok {
		return fmt.Errorf("invalid type for key 'ISAUTHENTICATED'")
	}
	isManagement, _ := c.Get(mgmnt_key).(bool)
	userName, _ := c.Get(username_key).(string)

	homeView := authviews.Usermanagement(isAuthenticated)
	c.Set("ISERROR", false)

	return render(c, authviews.UsermanagementIndex(
		"| Teams and Users",
		userName,
		isAuthenticated,
		isManagement,
		c.Get("ISERROR").(bool),
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		homeView,
	))
}

func (ah *AuthHandler) handleLoginTable(c echo.Context) error {
	isManagement, ok := c.Get(mgmnt_key).(bool)
	if !ok || !isManagement {
		return errors.New("Must be management!")
	}

	u := c.Get(user_id_key).(int32)
	users, _ := ah.UserServices.GetAllUsers()
	userTable := components.LoginTable(isManagement, users, u)

	return render(c, userTable)
}

func (ah *AuthHandler) homeHandler(c echo.Context) error {
	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	if !ok {
		return fmt.Errorf("invalid type for key 'ISAUTHENTICATED'")
	}
	isManagement, _ := c.Get(mgmnt_key).(bool)
	userName, _ := c.Get(username_key).(string)

	homeView := authviews.Home(isAuthenticated)
	c.Set("ISERROR", false)

	return render(c, authviews.HomeIndex(
		"| Home",
		userName,
		isAuthenticated,
		isManagement,
		c.Get("ISERROR").(bool),
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		homeView,
	))
}

func (ah *AuthHandler) flagsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(auth_sessions_key, c)

		if auth, ok := sess.Values[auth_key].(bool); !ok || !auth {
			fmt.Printf("\033[36m Ok=%t, Auth=%t \n\033[0m", ok, auth)
			c.Set("ISAUTHENTICATED", false)
			c.Set(mgmnt_key, false)

			return next(c)
		}

		if mgmt, ok := sess.Values[mgmnt_key].(bool); ok && mgmt {
			c.Set(mgmnt_key, true)
		}

		c.Set("ISAUTHENTICATED", true)

		return next(c)
	}
}

func (ah *AuthHandler) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(auth_sessions_key, c)
		if auth, ok := sess.Values[auth_key].(bool); !ok || !auth {
			// fmt.Printf("\033[36m Ok=%t, Auth=%t \n\033[0m", ok, auth)
			// fromProtected = false
			c.Set("ISAUTHENTICATED", false)

			return echo.NewHTTPError(echo.ErrUnauthorized.Code, "Please provide valid credentials")
		}

		if userId, ok := sess.Values[user_id_key].(int32); ok && userId != 0 {
			c.Set(user_id_key, userId) // set the user_id in the context
		}

		if username, ok := sess.Values[username_key].(string); ok && len(username) != 0 {
			c.Set(username_key, username) // set the username in the context
		}

		if tzone, ok := sess.Values[tzone_key].(string); ok && len(tzone) != 0 {
			c.Set(tzone_key, tzone) // set the client's time zone in the context
		}

		if isManagement, ok := sess.Values[mgmnt_key].(bool); ok && isManagement {
			c.Set(mgmnt_key, isManagement) // set the client's time zone in the context
		}
		// fromProtected = true
		c.Set("ISAUTHENTICATED", true)

		return next(c)
	}
}

func (ah *AuthHandler) loginHandler(c echo.Context) error {
	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	if !ok {
		return errors.New("invalid type for key 'ISAUTHENTICATED'")
	}
	loginView := authviews.Login(isAuthenticated)
	// isError = false
	c.Set("ISERROR", false)

	isManagement, ok := c.Get(mgmnt_key).(bool)
	if !ok {
		isManagement = false
	}
	if c.Request().Method == "POST" {
		// obtaining the time zone from the POST request of the login form
		tzone := ""
		if len(c.Request().Header["X-Timezone"]) != 0 {
			tzone = c.Request().Header["X-Timezone"][0]
			log.Printf("##################################################")
			log.Printf("TZone: %#v", tzone)
		}

		// Authentication goes here
		user, err := ah.UserServices.CheckUsername(c.FormValue("username"))
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				setFlashmessages(c, "error", "There is no user with that email")

				return c.Redirect(http.StatusSeeOther, "/login")
			}

			return echo.NewHTTPError(
				echo.ErrInternalServerError.Code,
				fmt.Sprintf(
					"something went wrong: %s",
					err,
				))
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(c.FormValue("password")),
		)
		if err != nil {
			// In production you have to give the user a generic message
			setFlashmessages(c, "error", "Incorrect password")

			return c.Redirect(http.StatusSeeOther, "/login")
		}

		// Get Session and setting Cookies
		sess, _ := session.Get(auth_sessions_key, c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   3600, // in seconds
			HttpOnly: true,
		}

		// Set user as authenticated, their username,
		// their ID and the client's time zone
		sess.Values = map[interface{}]interface{}{
			auth_key:     true,
			user_id_key:  user.ID,
			username_key: user.Name,
			tzone_key:    tzone,
			mgmnt_key:    user.IsManagement,
		}

		log.Printf("Session from loginHandler: %+v", sess)
		log.Printf("Session values: %+v", sess.Values)

		sess.Save(c.Request(), c.Response())

		setFlashmessages(c, "success", "You have successfully logged in!!")

		return c.Redirect(http.StatusSeeOther, "/")
	}

	return render(c, authviews.LoginIndex(
		"| Login",
		"",
		isAuthenticated,
		isManagement,
		c.Get("ISERROR").(bool),
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		loginView,
	))
}

func (ah *AuthHandler) logoutHandler(c echo.Context) error {
	sess, _ := session.Get(auth_sessions_key, c)
	// Revoke users authentication
	sess.Values = map[interface{}]interface{}{
		auth_key:     false,
		user_id_key:  "",
		username_key: "",
		tzone_key:    "",
		mgmnt_key:    false,
	}
	sess.Save(c.Request(), c.Response())

	setFlashmessages(c, "success", "You have successfully logged out!!")

	// fromProtected = false
	c.Set("ISAUTHENTICATED", false)

	return c.Redirect(http.StatusSeeOther, "/login")
}

func render(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
