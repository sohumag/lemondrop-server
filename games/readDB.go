package games

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	// defer func() {
	// 	if err := client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	return client
}

func (g *GameServer) GetAllCategories() ([]Category, error) {
	allSports, err := g.GetAllSportsFromDB()
	if err != nil {
		return nil, err
	}

	categories := make(map[string]bool)
	for _, sport := range allSports {
		if _, ok := categories[sport.Group]; !ok {
			categories[sport.Group] = true
		}
	}

	categoriesArr := []Category{}
	for k := range categories {
		categoriesArr = append(categoriesArr, Category{Name: k})
	}

	return categoriesArr, nil

}

func (g *GameServer) GetAllSportsInCategory(category string) ([]Sport, error) {
	allSports, err := g.GetAllSportsFromDB()
	if err != nil {
		return nil, err
	}

	matchedSports := []Sport{}
	for _, sport := range allSports {
		if sport.Group == category {
			matchedSports = append(matchedSports, sport)
		}
	}

	if len(matchedSports) == 0 {
		return nil, fmt.Errorf("Category does not exists")
	}

	return matchedSports, nil
}

func (g *GameServer) GetAllSportsFromDB() ([]Sport, error) {
	coll := g.client.Database("games-db").Collection("sports")
	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	var sports []Sport
	if err = cursor.All(context.TODO(), &sports); err != nil {
		log.Fatal(err)
	}

	fmt.Println(sports)
	return sports, nil
}

func (g *GameServer) GetUpcomingGamesBySport(sport string) ([]Game, error) {
	coll := g.client.Database("games-db").Collection("games")
	cursor, err := coll.Find(context.TODO(), bson.D{{"sporttitle", sport}})
	if err != nil {
		return nil, err
	}

	games := []Game{}
	if err = cursor.All(context.TODO(), &games); err != nil {
		return nil, err
	}

	if len(games) == 0 {
		return nil, fmt.Errorf("invalid sport")
	}

	return games, nil
}

func (g *GameServer) GetAllUpcomingGames(max int) ([]Game, error) {
	internalMaxNum := 25
	coll := g.client.Database("games-db").Collection("games")

	if max > internalMaxNum {
		max = internalMaxNum
	}

	opts := options.Find().SetLimit(int64(max))
	cursor, err := coll.Find(context.TODO(), bson.D{{}}, opts)
	if err != nil {
		return nil, err
	}

	games := []Game{}
	if err = cursor.All(context.TODO(), &games); err != nil {
		return nil, err
	}

	return games, nil

}
