package segment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type SourceMetadata struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Slug        string              `json:"slug"`
	Description string              `json:"description"`
	Logos       []Logo              `json:"logos"`
	Options     []IntegrationOption `json:"options"`
	Categories  []string            `json:"categories"`
}

type SourceSettings struct {
}

type Label struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type Source struct {
	ID          *string        `json:"id,omitempty"`
	Slug        string         `json:"slug"`
	Name        string         `json:"name"`
	Metadata    SourceMetadata `json:"metadata"`
	WorkspaceID string         `json:"workspaceId"`
	Enabled     bool           `json:"enabled"`
	WriteKeys   []string       `json:"writeKeys"`
	Settings    SourceSettings `json:"settings"`
	Labels      []Label        `json:"labels"`
}

type SourceResponse struct {
	Source Source `json:"source"`
}

func (c *Client) GetSource(sourceID string) (*Source, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/sources/%s", c.HostURL, sourceID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	sourceResponse := SourceResponse{}
	err = json.Unmarshal(body, &sourceResponse)
	if err != nil {
		return nil, err
	}

	return &sourceResponse.Source, nil
}

// func (c *Client) CreateSource(slug string, enabled bool, name string, metadataId string) (*Source, error) {
// 	newSource := Source{
//         Slug: slug,
//         Enabled: enabled,
//         Name: name,
//
// 	}
//
// 	newSourceData, err := json.Marshal(newSource)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sources/", c.HostURL, strings.NewReader(string(newSourceData)))
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	body, err := c.doRequest(req)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	connectionResponse := ConnectionResponse{}
// 	err = json.Unmarshal(body, &connectionResponse)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &connectionResponse.Data, nil
// }

// func (c *Client) UpdateConnection(connectionID, projectID string, connection Connection) (*Connection, error) {
// 	connectionData, err := json.Marshal(connection)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v3/accounts/%s/projects/%s/connections/%s/", c.HostURL, strconv.Itoa(c.AccountID), projectID, connectionID), strings.NewReader(string(connectionData)))
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	body, err := c.doRequest(req)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	connectionResponse := ConnectionResponse{}
// 	err = json.Unmarshal(body, &connectionResponse)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &connectionResponse.Data, nil
// }
//
// func (c *Client) DeleteConnection(connectionID, projectID string) (string, error) {
// 	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v3/accounts/%s/projects/%s/connections/%s/", c.HostURL, strconv.Itoa(c.AccountID), projectID, connectionID), nil)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	_, err = c.doRequest(req)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	return "", err
// }
