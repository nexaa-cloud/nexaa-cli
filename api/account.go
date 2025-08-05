package api

import (
	"strconv"

	"github.com/nexaa-cloud/nexaa-cli/config"
	"github.com/nexaa-cloud/nexaa-cli/graphql"
)

func GetAccountId() (int, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	var accountQuery struct {
		Account struct {
			Customer struct {
				Id string
			}
		}
	}

	params := map[string]graphql.Parameter{}

	query := client.BuildQuery(&accountQuery, params)

	err := client.Query(query)
	if err != nil {
		return 0, err
	}

	if id, err := strconv.Atoi(accountQuery.Account.Customer.Id); err == nil {
		return id, nil
	}

	return 0, err
}
