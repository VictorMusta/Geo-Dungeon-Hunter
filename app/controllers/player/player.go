package player

import (
	"dungeons/app/controllers/common"
	"dungeons/app/models"
	player "dungeons/app/services/player"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Player struct {
	PlayerService *player.Player
}

func New(playerService *player.Player) *Player {
	return &Player{
		PlayerService: playerService,
	}
}

// Get controller to get list of player
func (s *Player) Get(ctx *gin.Context) {
	var params models.QueryParams

	params.Parse(ctx)
	messageTypes := &models.MessageTypes{
		OK:                  "player.Search.Found",
		BadRequest:          "player.Search.BadRequest",
		NotFound:            "player.Search.NotFound",
		InternalServerError: "player.Search.Error",
	}

	players, err := s.PlayerService.Get(params)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}
	totalCount := len(players)
	if totalCount == 0 {
		status := http.StatusNotFound
		common.SendResponse(ctx, status, models.KnownError(status, messageTypes.NotFound, errors.New(" Data not found. ")))
		return
	}

	low := params.Offset - 1
	if low == -1 {
		low = 0
	}

	// Available CountMax calculation
	maxCount := params.Count
	if maxCount == 0 {
		maxCount = 100
	}

	high := maxCount + low
	if high > totalCount {
		high = totalCount
	}

	if low > high {
		status := http.StatusBadRequest
		common.SendResponse(ctx, status, models.KnownError(status, messageTypes.NotFound, errors.New(" Offset cannot be higher than count. ")))
		return
	}

	sendingPlayers := players[low:high]

	meta := models.MetaResponse{
		ObjectName: "Player",
		TotalCount: totalCount,
		Count:      len(sendingPlayers),
		Offset:     low + 1,
	}

	response := &models.WSResponse{
		Meta: meta,
		Data: sendingPlayers,
	}

	common.SendResponse(ctx, http.StatusOK, response)
}

// Create controller to create new player
func (s *Player) Create(ctx *gin.Context) {
	var in models.Player
	messageTypes := &models.MessageTypes{
		Created:             "player.Create.Created",
		BadRequest:          "player.Create.BadRequest",
		InternalServerError: "player.Create.Error",
	}

	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, messageTypes.BadRequest, err))
		return
	}

	player, err := s.PlayerService.Create(&in)
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}

	meta := models.MetaResponse{
		ObjectName: "Player",
		TotalCount: 1,
		Count:      1,
		Offset:     0,
	}
	response := &models.WSResponse{
		Meta: meta,
		Data: player,
	}

	common.SendResponse(ctx, http.StatusCreated, response)
}

// GetByID controller to get one player by id
func (s *Player) GetByID(ctx *gin.Context) {
	messageTypes := &models.MessageTypes{
		OK:                  "player.get.founded",
		NotFound:            "player.get.NotFound",
		BadRequest:          "player.get.BadRequest",
		InternalServerError: "player.get.Error",
	}

	id := ctx.Param("id")

	player, err := s.PlayerService.GetByID(id)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}

	// Créer la réponse avec le client
	meta := models.MetaResponse{
		ObjectName: "Player",
		TotalCount: 1,
		Count:      1,
		Offset:     0,
	}
	response := &models.WSResponse{
		Meta: meta,
		Data: player,
	}

	// Envoyer la réponse
	common.SendResponse(ctx, http.StatusOK, response)
}

// Update controller to update player
func (s *Player) Update(ctx *gin.Context) {
	var in models.Player
	messageTypes := &models.MessageTypes{
		OK:                  "player.Update.Updated",
		BadRequest:          "player.Update.BadRequest",
		InternalServerError: "players.Update.Error",
	}

	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, messageTypes.BadRequest, err))
		return
	}

	id := ctx.Param("id")
	err := s.PlayerService.Update(id, &in)
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}
	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, messageTypes.OK, "player updated"))
}

func (s *Player) Suspend(ctx *gin.Context) {
	messageTypes := &models.MessageTypes{
		OK:                  "player.Suspend.Updated",
		InternalServerError: "player.Suspend.Error",
	}
	id := ctx.Param("id")
	err := s.PlayerService.Suspend(id)
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}
	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, messageTypes.OK, "player suspended"))
}

// GetByIDs To get all player customID
func (s *Player) GetByIDs(ctx *gin.Context) {
	var params models.QueryParams
	params.Parse(ctx)
	messageTypes := &models.MessageTypes{
		OK:                  "player.Search.Updated",
		NotFound:            "player.Search.NotFound",
		BadRequest:          "player.Search.BadRequest",
		InternalServerError: "player.Search.Error",
	}

	// Extraction des identifiants de la requête URL
	ids := ctx.Param("ids")
	idList := strings.Split(ids, "&")

	players, err := s.PlayerService.GetByIds(idList)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}

	totalCount := len(players)
	if totalCount == 0 {
		status := http.StatusNotFound
		common.SendResponse(ctx, status, models.KnownError(status, messageTypes.NotFound, errors.New("Data not found.")))
		return
	}

	low := params.Offset - 1
	if low == -1 {
		low = 0
	}

	// Calcul du CountMax disponible
	maxCount := params.Count
	if maxCount == 0 {
		maxCount = 100
	}

	high := maxCount + low
	if high > totalCount {
		high = totalCount
	}

	if low > high {
		status := http.StatusBadRequest
		common.SendResponse(ctx, status, models.KnownError(status, messageTypes.NotFound, errors.New("Offset cannot be higher than count.")))
		return
	}

	sendingPlayer := players[low:high]

	meta := models.MetaResponse{
		ObjectName: "Player",
		TotalCount: totalCount,
		Count:      len(sendingPlayer),
		Offset:     low + 1,
	}

	response := &models.WSResponse{
		Meta: meta,
		Data: sendingPlayer,
	}

	common.SendResponse(ctx, http.StatusOK, response)
}
