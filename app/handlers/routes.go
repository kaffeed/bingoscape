package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func setupNonAuthRoutes(e *echo.Echo, ah *AuthHandler) {
	e.GET("/", ah.flagsMiddleware(ah.homeHandler))
	e.GET("/login", ah.flagsMiddleware(ah.loginHandler))
	e.POST("/login", ah.flagsMiddleware(ah.loginHandler))
	e.POST("/logout", ah.authMiddleware(ah.logoutHandler))
}

func setupLoginRoutes(tg *echo.Group, ah *AuthHandler) {
	tg.GET("", ah.handleUsermanagement)
	tg.GET("/create", ah.handleCreateLogin)
	tg.POST("/create", ah.handleCreateLogin)
	tg.GET("/list", ah.handleLoginTable)
	tg.DELETE("/:userId", ah.handleDeleteLogin)
	tg.PUT("/:userId/password", ah.handleChangePassword)
}

func setupTileRoutes(tg *echo.Group, th *TileHandler) {
	tg.GET("/:tileId", th.handleTile)
	tg.PUT("/:tileId", th.handleTile)
	tg.POST("/:tileId/submit", th.handleTileSubmission)
	tg.GET("/:tileId/templates", th.handleLoadFromTemplate)
	tg.GET("/:tileId/submissions", th.handleGetTileSubmissions)
	tg.PUT("/submissions/:submissionId/:state", th.handlePutSubmissionStatus)
	tg.DELETE("/:tileId/submissions/:submissionId", th.handleDeleteSubmission)
	tg.GET("/templates", th.handleGetTemplateTiles)
	tg.DELETE("/templates/:templateId", th.handleDeleteTemplate)
}

func setupBingoRoutes(bg *echo.Group, bh *BingoHandler) {
	bg.GET("/list", bh.handleGetAllBingos)
	bg.GET("/:bingoId", bh.handleGetBingoDetail)
	bg.Add("DIALOG", "/:bingoId", bh.handleGetBingoDetail)
	bg.GET("/create", bh.handleCreateBingo)
	bg.POST("/create", bh.handleCreateBingo)
	bg.DELETE("/delete/:bingoId", bh.handleDeleteBingo)
	bg.POST("/:bingoId/participants", bh.handleBingoParticipation)
	bg.GET("/:bingoId/participants", bh.handleBingoParticipation)
	bg.GET("/:bingoId/participantViewSwitch", bh.handleGetParticipationViewSwitch)
	bg.DELETE("/:bingoId/participants/:pId", bh.removeBingoParticipation)
	bg.PUT("/:bingoId/toggleState", bh.handleBingoState)
	bg.PUT("/:bingoId/toggleLeaderboardPublic", bh.handleBingoToggleLeaderboardPublic)
	bg.PUT("/:bingoId/toggleSubmissionsClosed", bh.handleBingoToggleSubmissionClosed)
	bg.GET("/:bingoId/teams/:loginId/submissions", bh.handleTeamSubmissions)
	bg.GET("/:bingoId/board", bh.handleGetBingoBoard)
}

func SetupRoutes(e *echo.Echo, authHandlers *AuthHandler, bingoHandlers *BingoHandler, tileHandler *TileHandler, apiHandler *ApiHandler) {
	setupNonAuthRoutes(e, authHandlers)

	teamGroup := e.Group("/logins", authHandlers.authMiddleware)
	setupLoginRoutes(teamGroup, authHandlers)

	tileGroup := e.Group("/tiles", authHandlers.authMiddleware)
	setupTileRoutes(tileGroup, tileHandler)

	// /* ↓ Protected Routes ↓ */
	bingoGroup := e.Group("/bingos", authHandlers.authMiddleware)
	setupBingoRoutes(bingoGroup, bingoHandlers)

	apiGroup := e.Group("/api", middleware.BasicAuth(basicAuthValidatorFunc(bingoHandlers.UserService)))

	setupApiRoutes(apiGroup, apiHandler)
}

func setupApiRoutes(ag *echo.Group, ah *ApiHandler) {
	ag.GET("/bingos/:bingoId", ah.handleGetBingo)
}
