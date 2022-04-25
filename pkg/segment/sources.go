package segment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type SourceMetadata struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Slug        string              `json:"slug"`
	Description string              `json:"description"`
	Logos       Logo                `json:"logos"`
	Options     []IntegrationOption `json:"options"`
	Categories  []string            `json:"categories"`
}

type TrackingSettings struct {
	AllowUnplannedEvents          bool   `json:"allow_unplanned_events"`
	AllowUnplannedEventProperties bool   `json:"allow_unplanned_event_properties"`
	AllowEventOnViolations        bool   `json:"allow_event_on_violations"`
	AllowPropertiesOnViolations   bool   `json:"allow_properties_on_violations"`
	CommonEventOnViolations       string `json:"common_event_on_violations"`
}

type IdentifySettings struct {
	AllowUnplannedTraits    bool   `json:"allow_unplanned_traits"`
	AllowTraitsOnViolations bool   `json:"allow_traits_on_violations"`
	CommonEventOnViolations string `json:"common_event_on_violations"`
}

type GroupSettings struct {
	AllowUnplannedTraits    bool   `json:"allow_unplanned_traits"`
	AllowTraitsOnViolations bool   `json:"allow_traits_on_violations"`
	CommonEventOnViolations string `json:"common_event_on_violations"`
}

type SourceSettings struct {
	ForwardingViolationsTo    string           `json:"forwarding_violations_to"`
	ForwardingBlockedEventsTo string           `json:"forwarding_blocked_events_to"`
	Track                     TrackingSettings `json:"track"`
	Identify                  IdentifySettings `json:"identify"`
	Group                     GroupSettings    `json:"group"`
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

type SourceResponseData struct {
	Data SourceResponse `json:"data"`
}

type SourceRequest struct {
	ID         *string        `json:"id,omitempty"`
	Slug       string         `json:"slug"`
	Name       string         `json:"name"`
	MetadataID *string        `json:"metadataId,omitempty"`
	Enabled    bool           `json:"enabled"`
	Settings   SourceSettings `json:"settings"`
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

	sourceResponseData := SourceResponseData{}
	err = json.Unmarshal(body, &sourceResponseData)
	if err != nil {
		return nil, err
	}

	return &sourceResponseData.Data.Source, nil
}

func (c *Client) CreateSource(slug string, enabled bool, name string, sourceSlug string, settings SourceSettings) (*Source, error) {
	sourceMetadata, _ := c.GetSourceMetadataFromCatalog(sourceSlug)

	newSource := SourceRequest{
		Slug:       slug,
		Enabled:    enabled,
		Name:       name,
		MetadataID: &sourceMetadata.ID,
		Settings:   settings,
	}

	newSourceData, err := json.Marshal(newSource)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sources/", c.HostURL), strings.NewReader(string(newSourceData)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	sourceResponseData := SourceResponseData{}
	err = json.Unmarshal(body, &sourceResponseData)
	if err != nil {
		return nil, err
	}

	return &sourceResponseData.Data.Source, nil
}

func (c *Client) UpdateSource(sourceID string, slug string, enabled bool, name string, settings SourceSettings) (*Source, error) {
	updatedSource := SourceRequest{
		ID:       &sourceID,
		Slug:     slug,
		Enabled:  enabled,
		Name:     name,
		Settings: settings,
	}

	updatedSourceData, err := json.Marshal(updatedSource)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/sources/%s", c.HostURL, sourceID), strings.NewReader(string(updatedSourceData)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	sourceResponseData := SourceResponseData{}
	err = json.Unmarshal(body, &sourceResponseData)
	if err != nil {
		return nil, err
	}

	return &sourceResponseData.Data.Source, nil
}

func (c *Client) DeleteSource(sourceID string) (string, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/sources/%s", c.HostURL, sourceID), nil)
	if err != nil {
		return "", err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return "", err
	}

	return "", err
}
