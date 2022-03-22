package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/kazakh-in-nz/m-bff/bff"
	"github.com/rs/zerolog/log"
)

func main() {

	grpcAddressHighScore := flag.String("address-m-highscore", "localhost:60051", "The grpc server address for highscore service")
	grpcAddressGameEngine := flag.String("address-m-game-engine", "localhost:60052", "The grpc server address for game engine service")

	serverAddress := flag.String("address-http", ":8081", "HTTP server address")

	flag.Parse()

	gameClient, errGameClient := bff.NewGrpcGameServiceClient(*grpcAddressHighScore)
	if errGameClient != nil {
		log.Error().Err(errGameClient).Msg("Error in creating a client for m-highscore")
	}

	gameEngineClient, errGameEngineClient := bff.NewGrpcGameEngineServiceClient(*grpcAddressGameEngine)
	if errGameEngineClient != nil {
		log.Error().Err(errGameEngineClient).Msg("Error in creating a client for m-game-engine")
	}

	gr := bff.NewGameResource(gameClient, gameEngineClient)

	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.GET("/geths", gr.GetHighScore)
	router.POST("/seths/:hs", gr.SetHighScore)
	router.GET("/getsize", gr.GetSize)
	router.POST("/setscore/:score", gr.SetScore)

	err := router.Run(*serverAddress)

	if err != nil {
		log.Fatal().Err(err).Msg("Could not start bff")
	}

	log.Info().Msgf("Started http-server at %v", *serverAddress)
}
