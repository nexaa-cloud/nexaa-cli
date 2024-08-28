package api

import (
	"net/http"

	"github.com/shurcooL/graphql"
	"gitlab.com/Tilaa/tilaa-cli/config"
)

var client *graphql.Client

func init() {
	httpClient := &http.Client{}
	client = graphql.NewClient(config.GRAPHQL_URL, httpClient)
}
