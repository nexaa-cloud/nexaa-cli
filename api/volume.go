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
	Locked bool
}

type VolumeInput struct {
	Id int
	Namespace string
	Name string
	Size int
}


type VolumeResponse struct {
    Id        string        `json:"id"`
    Name      string        `json:"name"`
    Namespace respNamespace `json:"namespace"`
    Size      int           `json:"size"`
    Usage     int           `json:"usage"`
	Locked	  bool			`json:"locked"`
}

type respNamespace struct {
    Name string `json:"name"`
}


func ListVolumes(namespace string) ([]Volume, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	var volumeQuery struct {
		Namespace struct {
			Id			string
			Name 		string
			Volumes []struct {
				Id 		string
				Name 	string
				Size 	int
				Usage 	int
				Locked	bool
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
			Size: int(volume.Size),
			Usage: int(volume.Usage),
			Locked: bool(volume.Locked),
		})
	}

	return volumes, nil
}

func ListVolumeByName(namespaceName string, volumeName string) (*Volume, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	var volumeQuery struct {
		Namespace struct {
			Name 	string
			Volumes []struct {
				Id		string
				Name 	string
				Size 	int
				Usage 	int
				Locked 	bool
			}
		} `graphql:"namespace(name: $name)"`
	}

	params := map[string]graphql.Parameter{
		"name": graphql.NewString(namespaceName),
	}

	query := client.BuildQuery(&volumeQuery, params)
	err := client.Query(query)

	if err != nil {
		return nil, err
	}

	var volume Volume
	
	for _, vol := range volumeQuery.Namespace.Volumes {
		if vol.Name == volumeName {
			volume.Id = vol.Id
			volume.Name = vol.Name
			volume.Size = vol.Size
			volume.Usage = vol.Usage
			volume.Locked = vol.Locked
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

	var resp VolumeResponse

	mutation := client.BuildMutationWithQuery("volumeCreate", params, &resp)

	err := client.Mutate(mutation)

	var vol Volume
	vol.Id = resp.Id
	vol.Name = resp.Name
	vol.Namespace = resp.Namespace.Name
	vol.Size = resp.Size
	vol.Usage = resp.Usage
	vol.Locked = resp.Locked

	return vol, err
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

	var resp VolumeResponse

	mutation := client.BuildMutationWithQuery("volumeIncrease", params, &resp)

	err := client.Mutate(mutation)

	var vol Volume
	vol.Id = resp.Id
	vol.Name = resp.Name
	vol.Namespace = resp.Namespace.Name
	vol.Size = resp.Size
	vol.Usage = resp.Usage
	vol.Locked = resp.Locked

	return vol, err
}


func DeleteVolume(name string, namespace string) error {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	volumeInput := map[string]any{
		"namespace": namespace,
		"name": name,
	}

	params := map[string]graphql.Parameter{
		"volume": graphql.NewComplexParameter("VolumeResourceInput", volumeInput),
	}

	mutation := client.BuildMutation("volumeDelete", params)

	err := client.Mutate(mutation)
	if err != nil {
		return err
	}

	return nil
}
