package handlers

import "github.com/labstack/echo/v4"

func SetupRoutes(e *echo.Echo, authHandlers *AuthHandler, bingoHandlers *BingoHandler, tileHandler *TileHandler) {
	e.GET("/", authHandlers.flagsMiddleware(authHandlers.homeHandler))
	e.GET("/login", authHandlers.flagsMiddleware(authHandlers.loginHandler))
	e.POST("/login", authHandlers.flagsMiddleware(authHandlers.loginHandler))
	e.POST("/logout", authHandlers.authMiddleware(authHandlers.logoutHandler))

	teamGroup := e.Group("/logins", authHandlers.authMiddleware)
	teamGroup.GET("", authHandlers.handleUsermanagement)
	teamGroup.GET("/create", authHandlers.handleCreateLogin)
	teamGroup.POST("/create", authHandlers.handleCreateLogin)
	teamGroup.GET("/list", authHandlers.handleLoginTable)
	teamGroup.DELETE("/:userId", authHandlers.handleDeleteLogin)
	teamGroup.PUT("/:userId/password", authHandlers.handleChangePassword)

	tileGroup := e.Group("/tiles", authHandlers.authMiddleware)
	tileGroup.GET("/:tileId", tileHandler.handleTile)
	tileGroup.PUT("/:tileId", tileHandler.handleTile)
	tileGroup.POST("/:tileId/submit", tileHandler.handleTileSubmission)
	tileGroup.GET("/:tileId/templates", tileHandler.handleLoadFromTemplate)
	tileGroup.GET("/:tileId/submissions", tileHandler.handleGetTileSubmissions)
	tileGroup.PUT("/submissions/:submissionId/:state", tileHandler.handlePutSubmissionStatus)
	tileGroup.DELETE("/:tileId/submissions/:submissionId", tileHandler.handleDeleteSubmission)
	tileGroup.GET("/templates", tileHandler.handleGetTemplateTiles)
	tileGroup.DELETE("/templates/:templateId", tileHandler.handleDeleteTemplate)

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
	protectedGroup.PUT("/:bingoId/toggleState", bingoHandlers.handleBingoState)
	protectedGroup.PUT("/:bingoId/toggleLeaderboardPublic", bingoHandlers.handleBingoToggleLeaderboardPublic)
	protectedGroup.PUT("/:bingoId/toggleSubmissionsClosed", bingoHandlers.handleBingoToggleSubmissionClosed)
	protectedGroup.GET("/:bingoId/teams/:loginId/submissions", bingoHandlers.handleTeamSubmissions)
	protectedGroup.GET("/:bingoId/board", bingoHandlers.handleGetBingoBoard)

	// protectedGroup.GET("/create", th.createTodoHandler)
	// protectedGroup.POST("/create", th.createTodoHandler)
	// protectedGroup.GET("/edit/:id", th.updateTodoHandler)
	// protectedGroup.POST("/edit/:id", th.updateTodoHandler)
	// protectedGroup.DELETE("/delete/:id", th.deleteTodoHandler)
}
