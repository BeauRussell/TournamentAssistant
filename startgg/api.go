package startgg

import (
	"log"

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

type Entrant struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type StandingsNode struct {
	Placement int     `json:"placement"`
	Entrant   Entrant `json:"entrant"`
}

type EventStandings struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Standings struct {
		Nodes []StandingsNode `json:"nodes"`
	} `json:"standings"`
}

func (s *Start) Setup(tournamentSlug string, key string) {
	s.client = graphql.NewClient("https://api.start.gg/gql/alpha", key)
	s.tournamentSlug = tournamentSlug
	s.apiKey = key
}

func (s *Start) GetEventData() *Tournament {
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

	return &respData.Data.Tournament
}

func (s *Start) GetEventStandings(eventId int) EventStandings {
	request := graphql.Request{
		Query: `query EventStandings($eventId: ID!, $page: Int!, $perPage: Int!) {
			event(id: $eventId) {
				id
				name
				standings(query: {
					perPage: $perPage,
					page: $page
				}){
					nodes {
						placement
						entrant {
							id
							name
						}
					}
				}
			}
		}`,
		Variables: map[string]interface{}{
			"eventId": eventId,
			"page":    0,
			"perPage": 20,
		},
	}

	var respData struct {
		Data struct {
			EventStandings EventStandings `json:"event"`
		} `json:"data"`
	}

	if err := s.client.Send(request, &respData); err != nil {
		log.Println("Failed to get Event Standings:", err)
	}

	return respData.Data.EventStandings
}
