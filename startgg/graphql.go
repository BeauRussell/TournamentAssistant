package startgg

import "github.com/machinebox/graphql"

type Start struct {
	client         *graphql.Client
	apiKey         string
	tournamentSlug string
}

func (s *Start) Setup(tournamentSlug string, key string) {
	s.client = graphql.NewClient("https://api.start.gg/gql/alpha")
	s.tournamentSlug = tournamentSlug
	s.apiKey = key
}
