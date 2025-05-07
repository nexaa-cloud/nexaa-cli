package api

import (
	// "github.com/shurcooL/graphql"

	"gitlab.com/tilaa/tilaa-cli/config"
	"gitlab.com/tilaa/tilaa-cli/graphql"
)

type Volume struct {
	Id string
	Namespace string
	Name string
	Size int
	Usage int
}

type VolumeInput struct {
	Id int
	Namespace string
	Name string
	Size int
}

func ListVolume(namespace string) ([]Volume, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	var volumeQuery struct {
		Namespace struct {
			Id			string
			Name 		string
			Volumes []struct {
				Id 		string
				Name 	string
				Size 	int
			}
		} `graphql:"namespace(name: $name)"`
	}

	params := map[string]graphql.Parameter{
		"name": graphql.NewString(namespace),
	}

	query := client.BuildQuery(&volumeQuery, params)
	err := client.Query(query)

	if err != nil {
		return nil, err
	}

	var volumes []Volume

	for _, volume := range volumeQuery.Namespace.Volumes {
		volumes = append(volumes, Volume{
			Id: string(volume.Id),
			Name: string(volume.Name),
		})
	}

	return volumes, nil
}

func ListVolumeById(namespace string, id string) (*Volume, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	var volumeQuery struct {
		Namespace struct {
			Name 	string
			Volumes []struct {
				Id		string
				Name 	string
				Size 	int
				Usage 	int
			}
		} `graphql:"namespace(id: $id)"`
	}

	params := map[string]graphql.Parameter{
		"id": graphql.NewId(namespace),
	}

	query := client.BuildQuery(&volumeQuery, params)
	err := client.Query(query)

	if err != nil {
		return nil, err
	}

	var volume Volume
	
	for _, vol := range volumeQuery.Namespace.Volumes {
		if vol.Id == id {
			volume.Id = vol.Id
			volume.Name = vol.Name
			volume.Size = vol.Size
			volume.Usage = vol.Usage
		}
	}

	volume.Namespace = volumeQuery.Namespace.Name

	return &volume, nil
}

func CreateVolume(input VolumeInput) (Volume, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	createVolumeInput := map[string]any{
		"name": input.Name,
		"namespace": input.Namespace,
		"size": input.Size,
	}
		
	params := map[string]graphql.Parameter{
		"volumeInput": graphql.NewComplexParameter("VolumeCreateInput", createVolumeInput),
	}

	mutation := client.BuildMutation("volumeCreate", params)

	err := client.Mutate(mutation)

	return Volume{}, err
}

func IncreaseVolume(input VolumeInput) (Volume, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	volumeInput := map[string]any{
		"name": input.Name,
		"namespace": input.Namespace,
		"size": input.Size,
	}

	params := map[string]graphql.Parameter{
		"volumeInput": graphql.NewComplexParameter("VolumeModifyInput", volumeInput),
	}

	mutation := client.BuildMutation("volumeIncrease", params)

	err := client.Mutate(mutation)

	return Volume{}, err
}


func DeleteVolume(name string, namespace string) error {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	volumeInput := map[string]any{
		"namespace": namespace,
		"name": name,
	}

	params := map[string]graphql.Parameter{
		"volumeDelete": graphql.NewComplexParameter("VolumeResourceinput", volumeInput),
	}

	mutation := client.BuildMutation("volumeDelete", params)

	err := client.Mutate(mutation)
	if err != nil {
		return err
	}

	return nil
}
