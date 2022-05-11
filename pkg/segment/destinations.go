package segment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type DestinationMetadata struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Logos       Logo   `json:"logos"`
	// 	Options     []IntegrationOption `json:"options"`
	Categories []string `json:"categories"`
}

type Destination struct {
	ID       *string                `json:"id,omitempty"`
	Name     string                 `json:"name"`
	Metadata SourceMetadata         `json:"metadata"`
	Enabled  bool                   `json:"enabled"`
	SourceID string                 `json:"sourceId"`
	Settings map[string]interface{} `json:"settings"`
}

type DestinationResponse struct {
	Destination Destination `json:"destination"`
}

type DestinationResponseData struct {
	Data DestinationResponse `json:"data"`
}

type DestinationRequest struct {
	ID         *string                `json:"id,omitempty"`
	SourceID   string                 `json:"sourceId"`
	Name       string                 `json:"name"`
	MetadataID *string                `json:"metadataId,omitempty"`
	Enabled    bool                   `json:"enabled"`
	Settings   map[string]interface{} `json:"settings"`
}

func (c *Client) GetDestination(destinationID string) (*Destination, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/destinations/%s", c.HostURL, destinationID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	destinationResponseData := DestinationResponseData{}
	err = json.Unmarshal(body, &destinationResponseData)
	if err != nil {
		return nil, err
	}

	return &destinationResponseData.Data.Destination, nil
}

func (c *Client) CreateDestination(sourceID string, enabled bool, name string, destinationSlug string, settings map[string]interface{}) (*Destination, error) {
	destinationMetadata, _ := c.GetDestinationMetadataFromCatalog(destinationSlug)

	newDestination := DestinationRequest{
		Enabled:    enabled,
		Name:       name,
		MetadataID: &destinationMetadata.ID,
		Settings:   settings,
		SourceID:   sourceID,
	}

	newDestinationData, err := json.Marshal(newDestination)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/destinations/", c.HostURL), strings.NewReader(string(newDestinationData)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	destinationResponseData := DestinationResponseData{}
	err = json.Unmarshal(body, &destinationResponseData)
	if err != nil {
		return nil, err
	}

	return &destinationResponseData.Data.Destination, nil
}

func (c *Client) UpdateDestination(destinationID string, sourceID string, enabled bool, name string, settings map[string]interface{}) (*Destination, error) {
	updatedDestination := DestinationRequest{
		SourceID: sourceID,
		Enabled:  enabled,
		Name:     name,
		Settings: settings,
	}

	updatedDestinationData, err := json.Marshal(updatedDestination)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/destinations/%s", c.HostURL, destinationID), strings.NewReader(string(updatedDestinationData)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	destinationResponseData := DestinationResponseData{}
	err = json.Unmarshal(body, &destinationResponseData)
	if err != nil {
		return nil, err
	}

	return &destinationResponseData.Data.Destination, nil
}

func (c *Client) DeleteDestination(destinationID string) (string, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/destinations/%s", c.HostURL, destinationID), nil)
	if err != nil {
		return "", err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return "", err
	}

	return "", err
}
