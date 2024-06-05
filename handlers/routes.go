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
	teamGroup.GET("", authHandlers.handleUsermanagement)
	teamGroup.GET("/create", bingoHandlers.RegisterHandler)
	teamGroup.POST("/create", bingoHandlers.RegisterHandler)
	teamGroup.GET("/list", authHandlers.handleLoginTable)
	// teamGroup.DELETE("/:id", bingoHandlers.handleDeleteLogins)

	tileGroup := e.Group("/tiles", authHandlers.authMiddleware)
	tileGroup.GET("/:tileId", bingoHandlers.handleTile)
	tileGroup.PUT("/:tileId", bingoHandlers.handleTile)
	tileGroup.POST("/:tileId/submit", bingoHandlers.handleTileSubmission)
	tileGroup.GET("/:tileId/submissions", bingoHandlers.handleGetTileSubmissions)
	tileGroup.PUT("/submissions/:submissionId/:state", bingoHandlers.handlePutSubmissionStatus)

	// /* ↓ Protected Routes ↓ */
	protectedGroup := e.Group("/bingos", authHandlers.authMiddleware)
	protectedGroup.GET("/list", bingoHandlers.handleGetAllBingos)
	protectedGroup.GET("/:bingoId", bingoHandlers.handleGetBingoDetail)
	protectedGroup.GET("/create", bingoHandlers.handleCreateBingo)
	protectedGroup.POST("/create", bingoHandlers.handleCreateBingo)
	protectedGroup.DELETE("/delete/:bingoId", bingoHandlers.handleDeleteBingo)
	protectedGroup.POST("/:bingoId/participants", bingoHandlers.handleBingoParticipation)
	protectedGroup.GET("/:bingoId/participants", bingoHandlers.handleBingoParticipation)
	protectedGroup.DELETE("/:bingoId/participants/:pId", bingoHandlers.removeBingoParticipation)

	// protectedGroup.GET("/create", th.createTodoHandler)
	// protectedGroup.POST("/create", th.createTodoHandler)
	// protectedGroup.GET("/edit/:id", th.updateTodoHandler)
	// protectedGroup.POST("/edit/:id", th.updateTodoHandler)
	// protectedGroup.DELETE("/delete/:id", th.deleteTodoHandler)
}
