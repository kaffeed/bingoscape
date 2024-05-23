package handlers

import "github.com/labstack/echo/v4"

func SetupRoutes(e *echo.Echo, authHandlers *AuthHandler, bingoHandlers *BingoHandler) {
	e.GET("/", authHandlers.flagsMiddleware(authHandlers.homeHandler))
	e.GET("/login", authHandlers.flagsMiddleware(authHandlers.loginHandler))
	e.POST("/login", authHandlers.flagsMiddleware(authHandlers.loginHandler))
	e.POST("/logout", authHandlers.authMiddleware(authHandlers.logoutHandler))
	// e.GET("/register", ah.flagsMiddleware(ah.registerHandler))
	// e.POST("/register", ah.flagsMiddleware(ah.registerHandler))

	// protectedGroup := e.Group("/bingo", authHandlers.authMiddleware)
	// /* ↓ Protected Routes ↓ */
	// protectedGroup.GET("/list", th.todoListHandler)
	// protectedGroup.GET("/create", th.createTodoHandler)
	// protectedGroup.POST("/create", th.createTodoHandler)
	// protectedGroup.GET("/edit/:id", th.updateTodoHandler)
	// protectedGroup.POST("/edit/:id", th.updateTodoHandler)
	// protectedGroup.DELETE("/delete/:id", th.deleteTodoHandler)
}
