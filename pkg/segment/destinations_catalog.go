package segment

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DestinationsCatalogResponseData struct {
	DestinationsCatalog []DestinationMetadata `json:"destinationsCatalog"`
	Pagination          Pagination            `json:"pagination"`
}

type DestinationsCatalogResponse struct {
	Data DestinationsCatalogResponseData `json:"data"`
}

func (c *Client) GetDestinationMetadataFromCatalog(destinationSlug string) (*DestinationMetadata, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/catalog/destinations?pagination.count=100", c.HostURL), nil)

	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)

	if err != nil {
		return nil, err
	}

	destinationsCatalogResponse := DestinationsCatalogResponse{}
	err = json.Unmarshal(body, &destinationsCatalogResponse)
	if err != nil {
		return nil, err
	}

	for i, destinationMetadata := range destinationsCatalogResponse.Data.DestinationsCatalog {
		if destinationMetadata.Slug == destinationSlug {
			return &destinationsCatalogResponse.Data.DestinationsCatalog[i], nil
		}
	}

	for {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/catalog/destinations?pagination.count=100&pagination.cursor=%s", c.HostURL, *destinationsCatalogResponse.Data.Pagination.Next), nil)
		if err != nil {
			return nil, err
		}

		body, err = c.doRequest(req)
		if err != nil {
			return nil, err
		}

		destinationsCatalogResponse = DestinationsCatalogResponse{}
		err = json.Unmarshal(body, &destinationsCatalogResponse)
		if err != nil {
			return nil, err
		}

		for i, destinationMetadata := range destinationsCatalogResponse.Data.DestinationsCatalog {
			if destinationMetadata.Slug == destinationSlug {
				return &destinationsCatalogResponse.Data.DestinationsCatalog[i], nil
			}
		}

		if *destinationsCatalogResponse.Data.Pagination.Next == "" {
			break
		}
	}

	return nil, fmt.Errorf("Did not find destination %s", destinationSlug)
}
