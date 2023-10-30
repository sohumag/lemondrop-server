package games

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *GameServer) ReturnPropsByIdAndPropName(c *fiber.Ctx, league string, gameId string, propName string) error {
	// needs to fetch if:
	// 		- document doesnt exist
	// 		- document exists but market does not

	// else: return market from document

	coll := s.client.Database("backup").Collection("props")
	filter := bson.D{{Key: "game_id", Value: gameId}}
	game := Game{}

	result := coll.FindOne(context.TODO(), filter)
	err := result.Decode(&game)

	if err == mongo.ErrNoDocuments {
		// create document with game id
		s.CreatePropDocumentById(gameId)
		market, err := s.GetPropsFromApi(league, gameId, propName)
		if err != nil {
			return err
		}
		s.AddPropDataToDB(gameId, propName, market, game.PropMarkets)
		c.JSON(market)
		return nil
	}

	// iterate through prop markets array and find if name matches: send back if does
	for _, propMarket := range game.PropMarkets {
		if propMarket.Key == propName {
			c.JSON(propMarket)
			return nil
		}
	}

	// if nothing mathces then it doesnt exist
	if err := s.AddPropDataToDBAndReturn(c, league, gameId, propName, game.PropMarkets); err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func (s *GameServer) CreatePropDocumentById(gameId string) error {
	coll := s.client.Database("backup").Collection("props")
	newGame := Game{Id: gameId}
	if _, err := coll.InsertOne(context.TODO(), newGame); err != nil {
		return err
	}

	return nil
}

func (s *GameServer) AddPropDataToDBAndReturn(c *fiber.Ctx, league string, gameId string, propName string, currentMarkets []Market) error {
	// expected document is already created
	propMarket, err := s.GetPropsFromApi(league, gameId, propName)
	if err != nil {
		{
			return err
		}
	}

	// add field to document
	s.AddPropDataToDB(gameId, propName, propMarket, currentMarkets)
	c.JSON(propMarket)
	return nil
}

func (s *GameServer) AddPropDataToDB(gameId string, propName string, market Market, allMarkets []Market) error {
	coll := s.client.Database("backup").Collection("props")
	newMarkets := append(allMarkets, market)
	//!!! key component
	if _, err := coll.UpdateOne(context.TODO(), bson.D{{Key: "game_id", Value: gameId}}, bson.D{{Key: "$set", Value: bson.D{{Key: "prop_markets", Value: newMarkets}}}}); err != nil {
		return err
	}

	return nil
}

func (s *GameServer) GetPropsFromApi(league, gameId, propName string) (Market, error) {
	apiKey := os.Getenv("ODDS_API_KEY")
	url := fmt.Sprintf("https://api.the-odds-api.com/v4/sports/%s/events/%s/odds?apiKey=%s&regions=us&markets=%s&oddsFormat=american", league, gameId, apiKey, propName)
	res, err := http.Get(url)
	if err != nil {
		return Market{}, err
	}
	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return Market{}, err
	}

	game := Game{}
	json.Unmarshal(bytes, &game)
	propMarket := game.Bookmakers[0].Markets[0]
	return propMarket, nil
}
