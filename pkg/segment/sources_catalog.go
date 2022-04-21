package segment

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SourcesCatalogResponseData struct {
	SourcesCatalog []SourceMetadata `json:"sourcesCatalog"`
	Pagination     Pagination       `json:"pagination"`
}

type SourcesCatalogResponse struct {
	Data SourcesCatalogResponseData `json:"default"`
}

func (c *Client) GetSourceMetadataFromCatalog(sourceSlug string) (*SourceMetadata, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/sources/catalog?pagination.count=100", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	sourcesCatalogResponse := SourcesCatalogResponse{}
	err = json.Unmarshal(body, &sourcesCatalogResponse)
	if err != nil {
		return nil, err
	}
	for i, sourceMetadata := range sourcesCatalogResponse.Data.SourcesCatalog {
		if sourceMetadata.Slug == sourceSlug {
			return &sourcesCatalogResponse.Data.SourcesCatalog[i], nil
		}
	}

	for {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/sources/catalog?pagination.count=100&pagination.cursor=%s", c.HostURL, sourcesCatalogResponse.Data.Pagination.Next), nil)
		if err != nil {
			return nil, err
		}
		body, err = c.doRequest(req)
		if err != nil {
			return nil, err
		}
		sourcesCatalogResponse = SourcesCatalogResponse{}
		err = json.Unmarshal(body, &sourcesCatalogResponse)
		if err != nil {
			return nil, err
		}
		for i, sourceMetadata := range sourcesCatalogResponse.Data.SourcesCatalog {
			if sourceMetadata.Slug == sourceSlug {
				return &sourcesCatalogResponse.Data.SourcesCatalog[i], nil
			}
		}

		if *sourcesCatalogResponse.Data.Pagination.Next == "" {
			break
		}
	}

	return nil, fmt.Errorf("Did not find source %s", sourceSlug)
}
