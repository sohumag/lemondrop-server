package games

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type GameServer struct {
	port   string
	client *mongo.Client
}

func NewGameServer(listenAddr string) *GameServer {
	// g := &GameServer{port: listenAddr}

	return &GameServer{
		port:   listenAddr,
		client: ConnectDB(),
	}
}

func (g *GameServer) Start() error {
	g.StartAPI()
	return nil
}

type Category struct {
	Name string `json:"name"`
}

type Sport struct {
	Key         string `json:"key"`
	Title       string `json:"title"`
	Group       string `json:"group"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type Game struct {
	GameId       string      `json:"id"`
	SportKey     string      `json:"sport_key"`
	SportTitle   string      `json:"sport_title"`
	CommenceTime time.Time   `json:"commence_time"`
	HomeTeam     string      `json:"home_team"`
	AwayTeam     string      `json:"away_team"`
	Bookmakers   []Bookmaker `json:"bookmakers"`
}

type Bookmaker struct {
	Key        string    `json:"key"`
	Title      string    `json:"title"`
	LastUpdate time.Time `json:"last_update"`
	Markets    []Market  `json:"markets"`
}
type Market struct {
	Key        string    `json:"key"`
	LastUpdate time.Time `json:"last_update"`
	Outcomes   []Outcome `json:"outcomes"`
}

type Outcome struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Point float64 `json:"point"`
}

// func (g *GameServer) GetAllSports() ([]Sport, error) {
// 	path := "./allsports.json"
// 	bytes, err := os.ReadFile(path)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var sports []Sport
// 	if err = json.Unmarshal([]byte(bytes), &sports); err != nil {
// 		return nil, err
// 	}

// 	return sports, nil
// }

// func (g *GameServer) GetSportsCategories() []Category {
// 	allSports, _ := g.GetAllSports()

// 	categories := make(map[string]bool)
// 	for _, sport := range allSports {
// 		if _, ok := categories[sport.Group]; !ok {
// 			categories[sport.Group] = true
// 		}
// 	}

// 	categoriesArr := []Category{}
// 	for k := range categories {
// 		categoriesArr = append(categoriesArr, Category{Name: k})
// 	}

// 	return categoriesArr
// }
