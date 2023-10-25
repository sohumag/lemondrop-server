package archivedGames

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

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

func (g *GameServer) AddSportToDB(sport *Sport) error {
	coll := g.client.Database("games-db").Collection("sports")
	result, err := coll.InsertOne(context.TODO(), *sport)
	if err != nil {
		return err
	}

	fmt.Printf("Added sport with id: %v\n", result.InsertedID)

	return nil
}
