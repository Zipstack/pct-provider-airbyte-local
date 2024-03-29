package api

import (
	"encoding/json"
	"fmt"
)

type SourceStripeID struct {
	SourceId string `json:"sourceId"`
}

type SourceStripe struct {
	Name                    string                 `json:"name"`
	SourceId                string                 `json:"sourceId,omitempty"`
	SourceDefinitionId      string                 `json:"sourceDefinitionId,omitempty"`
	WorkspaceId             string                 `json:"workspaceId,omitempty"`
	ConnectionConfiguration SourceStripeConnConfig `json:"connectionConfiguration"`
}

type SourceStripeConnConfig struct {
	StartDate          string `json:"start_date"`
	LookbackWindowDays int    `json:"lookback_window_days,omitempty"`
	SliceRange         int    `json:"slice_range,omitempty"`
	ClientSecret       string `json:"client_secret"`
	AccountId          string `json:"account_id"`
}

func (c *Client) CreateStripeSource(payload SourceStripe) (SourceStripe, error) {
	// logger := fwhelpers.GetLogger()
	method := "POST"
	url := c.Host + "/api/v1/sources/create"
	body, err := json.Marshal(payload)
	if err != nil {
		return SourceStripe{}, err
	}

	b, statusCode, _, _, err := c.doRequest(method, url, body, nil)
	if err != nil {
		return SourceStripe{}, err
	}
	source := SourceStripe{}
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

func (c *Client) ReadStripeSource(sourceId string) (SourceStripe, error) {
	// logger := fwhelpers.GetLogger()

	method := "POST"
	url := c.Host + "/api/v1/sources/get"
	sId := SourceStripeID{sourceId}
	body, err := json.Marshal(sId)
	if err != nil {
		return SourceStripe{}, err
	}

	b, statusCode, _, _, err := c.doRequest(method, url, body, nil)
	if err != nil {
		return SourceStripe{}, err
	}

	source := SourceStripe{}
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

func (c *Client) UpdateStripeSource(payload SourceStripe) (SourceStripe, error) {
	// logger := fwhelpers.GetLogger()

	method := "POST"
	url := c.Host + "/api/v1/sources/update"
	body, err := json.Marshal(payload)
	if err != nil {
		return SourceStripe{}, err
	}

	b, statusCode, _, _, err := c.doRequest(method, url, body, nil)
	if err != nil {
		return SourceStripe{}, err
	}

	source := SourceStripe{}
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

func (c *Client) DeleteStripeSource(sourceId string) error {
	// logger := fwhelpers.GetLogger()

	method := "POST"
	url := c.Host + "/api/v1/sources/delete"
	sId := SourceStripeID{sourceId}
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
