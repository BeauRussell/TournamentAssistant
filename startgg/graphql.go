package startgg

import "github.com/machinebox/graphql"

type Start struct {
	client *graphql.Client
	apiKey string
}

func (s *Start) Setup() {
	s.client = graphql.NewClient("https://api.start.gg/gql/alpha")
}
