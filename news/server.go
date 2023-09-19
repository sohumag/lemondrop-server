package news

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type NewsServer struct{}

func NewNewsServer() *NewsServer {
	return &NewsServer{}
}

func (s *NewsServer) Start(api fiber.Router) error {
	// s.StartNewsServerAPI(api)
	s.GetNews()
	return nil
}

func (s *NewsServer) GetNews() error {
	newsEndpoints := map[string]string{
		"College Football": "https://sports.yahoo.com/college-football/news/",
		"NBA":              "https://sports.yahoo.com/nba/news/",
		"MLB":              "https://sports.yahoo.com/mlb/news/",
	}

	wg := sync.WaitGroup{}
	wg.Add(len(newsEndpoints))

	for category, url := range newsEndpoints {
		go func(category string, url string) {
			res, err := http.Get(url)
			if err != nil {
				panic(err)
			}

			defer res.Body.Close()

			bytes, _ := io.ReadAll(res.Body)

			if err = os.WriteFile(fmt.Sprintf("./news/%s.html", category), bytes, 0644); err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}(category, url)
	}

	return nil
}
