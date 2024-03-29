package api

import (
	"encoding/json"
	"fmt"
)

type SourcePipedriveID struct {
	SourceId string `json:"sourceId"`
}

type SourcePipedrive struct {
	Name                    string                    `json:"name"`
	SourceId                string                    `json:"sourceId,omitempty"`
	SourceDefinitionId      string                    `json:"sourceDefinitionId,omitempty"`
	WorkspaceId             string                    `json:"workspaceId,omitempty"`
	ConnectionConfiguration SourcePipedriveConnConfig `json:"connectionConfiguration"`
}

type SourcePipedriveConnConfig struct {
	ReplicationStartDate string                         `json:"replication_start_date"`
	Authorization        SourcePipedriveAuthConfigModel `json:"authorization"`
}

type SourcePipedriveAuthConfigModel struct {
	AuthType string `json:"auth_type"`
	ApiToken string `json:"api_token"`
}

func (c *Client) CreatePipedriveSource(payload SourcePipedrive) (SourcePipedrive, error) {
	// logger := fwhelpers.GetLogger()

	method := "POST"
	url := c.Host + "/api/v1/sources/create"
	body, err := json.Marshal(payload)
	if err != nil {
		return SourcePipedrive{}, err
	}

	b, statusCode, _, _, err := c.doRequest(method, url, body, nil)
	if err != nil {
		return SourcePipedrive{}, err
	}

	source := SourcePipedrive{}
	if statusCode >= 200 && statusCode <= 299 {
		err = json.Unmarshal(b, &source)
		return source, err
	} else {
		msg, err := c.getAPIError(b)
		if err != nil {
			return source, err
		} else {
			return source, fmt.Errorf(msg)
		}
	}
}

func (c *Client) ReadPipedriveSource(sourceId string) (SourcePipedrive, error) {
	// logger := fwhelpers.GetLogger()

	method := "POST"
	url := c.Host + "/api/v1/sources/get"
	sId := SourcePipedriveID{sourceId}
	body, err := json.Marshal(sId)
	if err != nil {
		return SourcePipedrive{}, err
	}

	b, statusCode, _, _, err := c.doRequest(method, url, body, nil)
	if err != nil {
		return SourcePipedrive{}, err
	}

	source := SourcePipedrive{}
	if statusCode >= 200 && statusCode <= 299 {
		err = json.Unmarshal(b, &source)
		return source, err
	} else {
		msg, err := c.getAPIError(b)
		if err != nil {
			return source, err
		} else {
			return source, fmt.Errorf(msg)
		}
	}
}

func (c *Client) UpdatePipedriveSource(payload SourcePipedrive) (SourcePipedrive, error) {
	// logger := fwhelpers.GetLogger()

	method := "POST"
	url := c.Host + "/api/v1/sources/update"
	body, err := json.Marshal(payload)
	if err != nil {
		return SourcePipedrive{}, err
	}

	b, statusCode, _, _, err := c.doRequest(method, url, body, nil)
	if err != nil {
		return SourcePipedrive{}, err
	}

	source := SourcePipedrive{}
	if statusCode >= 200 && statusCode <= 299 {
		err = json.Unmarshal(b, &source)
		return source, err
	} else {
		msg, err := c.getAPIError(b)
		if err != nil {
			return source, err
		} else {
			return source, fmt.Errorf(msg)
		}
	}
}

func (c *Client) DeletePipedriveSource(sourceId string) error {
	// logger := fwhelpers.GetLogger()

	method := "POST"
	url := c.Host + "/api/v1/sources/delete"
	sId := SourcePipedriveID{sourceId}
	body, err := json.Marshal(sId)
	if err != nil {
		return err
	}

	b, statusCode, _, _, err := c.doRequest(method, url, body, nil)
	if err != nil {
		return err
	}

	if statusCode >= 200 && statusCode <= 299 {
		return nil
	} else {
		msg, err := c.getAPIError(b)
		if err != nil {
			return err
		} else {
			return fmt.Errorf(msg)
		}
	}
}
