package bff

import (
	"context"
	"fmt"
	"strconv"
	"time"

	pbgameengine "github.com/kazakh-in-nz/m-game-engine/v1"
	pbhighscore "github.com/kazakh-in-nz/m_apis/v1"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type gameResource struct {
	gameClient       pbhighscore.GameClient
	gameEngineClient pbgameengine.GameEngineClient
}

func NewGameResource(gameClient pbhighscore.GameClient, gameEngineClient pbgameengine.GameEngineClient) *gameResource {
	return &gameResource{
		gameClient:       gameClient,
		gameEngineClient: gameEngineClient,
	}
}

func NewGrpcGameServiceClient(serverAddr string) (pbhighscore.GameClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())

	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("Failed to dial: %v", err))
		return nil, err
	} else {
		log.Info().Msgf("Successfully connected to [%s]", serverAddr)
	}

	if conn == nil {
		log.Info().Msg("m-highscore connection is nil in m-bff")
	}

	client := pbhighscore.NewGameClient(conn)

	return client, nil
}

func NewGrpcGameEngineServiceClient(serverAddr string) (pbgameengine.GameEngineClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())

	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("Failed to dial: %v", err))
		return nil, err
	} else {
		log.Info().Msgf("Successfully connected to [%s]", serverAddr)
	}

	if conn == nil {
		log.Info().Msg("m-game-engine connection is nil in m-bff")
	}

	client := pbgameengine.NewGameEngineClient(conn)

	return client, nil
}

func (gr *gameResource) GetHighScore(c *gin.Context) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	highScoreResponse, err := gr.gameClient.GetHighScore(timeoutCtx, &pbhighscore.GetHighScoreRequest{})

	if err != nil {
		log.Info().Msg("Test message ======>")
		log.Error().Err(err).Msg("Failed to get high score")
		log.Panic()
	}

	hsString := strconv.FormatFloat(highScoreResponse.GetHighScore(), 'e', -1, 64)

	c.JSONP(200, gin.H{
		"hs": hsString,
	})
}

func (gr *gameResource) SetHighScore(c *gin.Context) {
	highScoreString := c.Param("hs")
	highScoreFloat64, err := strconv.ParseFloat(highScoreString, 64)

	if err != nil {
		log.Error().Err(err).Msg("Failed to convert highscore to float")
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, error := gr.gameClient.SetHighScore(timeoutCtx, &pbhighscore.SetHighScoreRequest{
		HighScore: highScoreFloat64,
	})

	if error != nil {
		log.Error().Err(err).Msg("Error while setting high score in m-highscore")
	}
}

func (gr *gameResource) GetSize(c *gin.Context) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	sizeResponse, err := gr.gameEngineClient.GetSize(timeoutCtx, &pbgameengine.GetSizeRequest{})

	if err != nil {
		log.Error().Err(err).Msg("Failed to get size")
		log.Panic()
	}

	sizeString := strconv.FormatFloat(sizeResponse.GetSize(), 'e', -1, 64)

	c.JSONP(200, gin.H{
		"size": sizeString,
	})
}

func (gr *gameResource) SetScore(c *gin.Context) {
	scoreString := c.Param("score")
	score64, err := strconv.ParseFloat(scoreString, 64)

	if err != nil {
		log.Error().Err(err).Msg("Failed to convert score to float")
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, error := gr.gameEngineClient.SetScore(timeoutCtx, &pbgameengine.SetScoreRequest{
		Score: score64,
	})

	if error != nil {
		log.Error().Err(err).Msg("Error while setting score in m-game-engine")
	}
}
