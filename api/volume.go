package api

import (
	"gitlab.com/tilaa/tilaa-cli/config"
	"gitlab.com/tilaa/tilaa-cli/graphql"
)

type Volume struct {
	Namespace string
	Name string
	Size int
	Usage int
	Locked bool
}

type VolumeInput struct {
	Namespace string
	Name string
	Size int
}

type VolumeResponse struct {
    Name      string        	 `json:"name"`
    Namespace NamespaceResponse	 `json:"namespace"`
    Size      int           	 `json:"size"`
    Usage     int           	 `json:"usage"`
	Locked	  bool			 	 `json:"locked"`
}


func ListVolumes(namespace string) ([]Volume, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	var volumeQuery struct {
		Namespace struct {
			Name 		string
			Volumes []struct {
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
			Name: string(volume.Name),
			Size: int(volume.Size),
			Usage: int(volume.Usage),
			Locked: bool(volume.Locked),
		})
	}

	return volumes, nil
}

func ListVolumeByName(namespace string, volume string) (*Volume, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	var volumeQuery struct {
		Namespace struct {
			Name 	string
			Volumes []struct {
				Name 	string
				Size 	int
				Usage 	int
				Locked 	bool
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

	var vol Volume
	
	for _, item := range volumeQuery.Namespace.Volumes {
		if item.Name == volume {
			vol.Name = item.Name
			vol.Size = item.Size
			vol.Usage = item.Usage
			vol.Locked = item.Locked
		}
	}

	vol.Namespace = volumeQuery.Namespace.Name

	return &vol, nil
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
