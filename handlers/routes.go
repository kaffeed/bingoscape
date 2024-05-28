package handlers

import "github.com/labstack/echo/v4"

func SetupRoutes(e *echo.Echo, authHandlers *AuthHandler, bingoHandlers *BingoHandler) {
	e.GET("/", authHandlers.flagsMiddleware(authHandlers.homeHandler))
	e.GET("/login", authHandlers.flagsMiddleware(authHandlers.loginHandler))
	e.POST("/login", authHandlers.flagsMiddleware(authHandlers.loginHandler))
	e.POST("/logout", authHandlers.authMiddleware(authHandlers.logoutHandler))
	// e.GET("/register", ah.flagsMiddleware(ah.registerHandler))
	// e.POST("/register", ah.flagsMiddleware(ah.registerHandler))

	teamGroup := e.Group("/logins", authHandlers.authMiddleware)
	teamGroup.GET("/create", bingoHandlers.RegisterHandler)
	teamGroup.POST("/create", bingoHandlers.RegisterHandler)

	tileGroup := e.Group("/tiles", authHandlers.authMiddleware)
	tileGroup.GET("/:tileId/edit", bingoHandlers.handleEditTile)
	tileGroup.PUT("/:tileId/edit", bingoHandlers.handleEditTile)

	// /* ↓ Protected Routes ↓ */
	protectedGroup := e.Group("/bingos", authHandlers.authMiddleware)
	protectedGroup.GET("/list", bingoHandlers.handleGetAllBingos)
	protectedGroup.GET("/:bingoId", bingoHandlers.handleGetBingoDetail)
	protectedGroup.GET("/create", bingoHandlers.handleCreateBingo)
	protectedGroup.POST("/create", bingoHandlers.handleCreateBingo)
	protectedGroup.DELETE("/delete/:bingoId", bingoHandlers.handleDeleteBingo)
	protectedGroup.GET("/:bingoId/active/:isActive", bingoHandlers.handleUpdateActiveState)
	protectedGroup.PUT("/:bingoId/active/:isActive", bingoHandlers.handleUpdateActiveState)
	protectedGroup.POST("/:bingoId/participants", bingoHandlers.handleBingoParticipation)
	protectedGroup.GET("/:bingoId/participants", bingoHandlers.handleBingoParticipation)
	protectedGroup.DELETE("/:bingoId/participants/:pId", bingoHandlers.removeBingoParticipation)

	// protectedGroup.GET("/create", th.createTodoHandler)
	// protectedGroup.POST("/create", th.createTodoHandler)
	// protectedGroup.GET("/edit/:id", th.updateTodoHandler)
	// protectedGroup.POST("/edit/:id", th.updateTodoHandler)
	// protectedGroup.DELETE("/delete/:id", th.deleteTodoHandler)
}
