package startgg

import (
	"log"
	"net/http"

	"github.com/BeauRussell/TournamentAssistant/graphql"
)

type Start struct {
	client         *graphql.Client
	apiKey         string
	tournamentSlug string
}

type Event struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Tournament struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Events []Event `json:"events"`
}

type loggingTransport struct {
	Transport http.RoundTripper
}

func (s *Start) Setup(tournamentSlug string, key string) {
	s.client = graphql.NewClient("https://api.start.gg/gql/alpha", key)
	s.tournamentSlug = tournamentSlug
	s.apiKey = key
}

func (s *Start) GetEventData() Tournament {
	request := graphql.Request{
		Query: `query TournamentEvents($tourneySlug: String!) {
			tournament(slug: $tourneySlug) {
				id
				name
				events {
					id
					name
				}
			}
		}`,
		Variables: map[string]interface{}{
			"tourneySlug": s.tournamentSlug,
		},
	}

	var respData struct {
		Data struct {
			Tournament Tournament `json:"tournament"`
		} `json:"data"`
	}

	if err := s.client.Send(request, &respData); err != nil {
		log.Println("Failed to get Event Data:", err)
		panic(err)
	}

	return respData.Data.Tournament
}
