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
	bingoGroup := e.Group("/bingos", authHandlers.authMiddleware)
	bingoGroup.GET("/list", bingoHandlers.handleGetAllBingos)
	bingoGroup.GET("/:bingoId", bingoHandlers.handleGetBingoDetail)
	bingoGroup.Add("DIALOG", "/:bingoId", bingoHandlers.handleGetBingoDetail)
	bingoGroup.GET("/create", bingoHandlers.handleCreateBingo)
	bingoGroup.POST("/create", bingoHandlers.handleCreateBingo)
	bingoGroup.DELETE("/delete/:bingoId", bingoHandlers.handleDeleteBingo)
	bingoGroup.POST("/:bingoId/participants", bingoHandlers.handleBingoParticipation)
	bingoGroup.GET("/:bingoId/participants", bingoHandlers.handleBingoParticipation)
	bingoGroup.GET("/:bingoId/participantViewSwitch", bingoHandlers.handleGetParticipationViewSwitch)
	bingoGroup.DELETE("/:bingoId/participants/:pId", bingoHandlers.removeBingoParticipation)
	bingoGroup.PUT("/:bingoId/toggleState", bingoHandlers.handleBingoState)
	bingoGroup.PUT("/:bingoId/toggleLeaderboardPublic", bingoHandlers.handleBingoToggleLeaderboardPublic)
	bingoGroup.PUT("/:bingoId/toggleSubmissionsClosed", bingoHandlers.handleBingoToggleSubmissionClosed)
	bingoGroup.GET("/:bingoId/teams/:loginId/submissions", bingoHandlers.handleTeamSubmissions)
	bingoGroup.GET("/:bingoId/board", bingoHandlers.handleGetBingoBoard)

	// protectedGroup.GET("/create", th.createTodoHandler)
	// protectedGroup.POST("/create", th.createTodoHandler)
	// protectedGroup.GET("/edit/:id", th.updateTodoHandler)
	// protectedGroup.POST("/edit/:id", th.updateTodoHandler)
	// protectedGroup.DELETE("/delete/:id", th.deleteTodoHandler)
}
