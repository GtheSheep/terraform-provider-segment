package segment

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WarehousesCatalogResponseData struct {
	WarehousesCatalog []WarehouseMetadata `json:"warehousesCatalog"`
	Pagination        Pagination          `json:"pagination"`
}

type WarehousesCatalogResponse struct {
	Data WarehousesCatalogResponseData `json:"data"`
}

func (c *Client) GetWarehouseMetadataFromCatalog(warehouseSlug string) (*WarehouseMetadata, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/catalog/warehouses?pagination.count=100", c.HostURL), nil)

	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)

	if err != nil {
		return nil, err
	}

	warehousesCatalogResponse := WarehousesCatalogResponse{}
	err = json.Unmarshal(body, &warehousesCatalogResponse)
	if err != nil {
		return nil, err
	}

	for i, warehouseMetadata := range warehousesCatalogResponse.Data.WarehousesCatalog {
		if warehouseMetadata.Slug == warehouseSlug {
			return &warehousesCatalogResponse.Data.WarehousesCatalog[i], nil
		}
	}

	for {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/warehouses/catalog?pagination.count=100&pagination.cursor=%s", c.HostURL, *warehousesCatalogResponse.Data.Pagination.Next), nil)
		if err != nil {
			return nil, err
		}

		body, err = c.doRequest(req)
		if err != nil {
			return nil, err
		}

		warehousesCatalogResponse = WarehousesCatalogResponse{}
		err = json.Unmarshal(body, &warehousesCatalogResponse)
		if err != nil {
			return nil, err
		}

		for i, warehouseMetadata := range warehousesCatalogResponse.Data.WarehousesCatalog {
			if warehouseMetadata.Slug == warehouseSlug {
				return &warehousesCatalogResponse.Data.WarehousesCatalog[i], nil
			}
		}

		if *warehousesCatalogResponse.Data.Pagination.Next == "" {
			break
		}
	}

	return nil, fmt.Errorf("Did not find warehouse %s", warehouseSlug)
}
