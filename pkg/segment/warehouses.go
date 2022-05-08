package segment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type WarehouseMetadata struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Slug        string              `json:"slug"`
	Description string              `json:"description"`
	Logos       Logo                `json:"logos"`
	Options     []IntegrationOption `json:"options"`
}

type WarehouseSettings struct {
	Hostname   string `json:"hostname,omitempty"`
	Database   string `json:"database,omitempty"`
	Port       string `json:"port,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	Ciphertext string `json:"ciphertext,omitempty"`
	Name       string `json:"name,omitempty"`
}

type Warehouse struct {
	ID          *string           `json:"id,omitempty"`
	Metadata    WarehouseMetadata `json:"metadata"`
	Name        string            `json:"name"`
	WorkspaceID string            `json:"workspaceId"`
	Enabled     bool              `json:"enabled"`
	Settings    WarehouseSettings `json:"settings"`
}

type WarehouseResponse struct {
	Warehouse Warehouse `json:"warehouse"`
}

type WarehouseResponseData struct {
	Data WarehouseResponse `json:"data"`
}

type WarehouseRequest struct {
	ID         *string           `json:"id,omitempty"`
	Name       string            `json:"name"`
	MetadataID *string           `json:"metadataId,omitempty"`
	Enabled    bool              `json:"enabled"`
	Settings   WarehouseSettings `json:"settings"`
}

func (c *Client) GetWarehouse(warehouseID string) (*Warehouse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/warehouses/%s", c.HostURL, warehouseID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	warehouseResponseData := WarehouseResponseData{}
	err = json.Unmarshal(body, &warehouseResponseData)
	if err != nil {
		return nil, err
	}

	if warehouseResponseData.Data.Warehouse.Name == "" {
		warehouseResponseData.Data.Warehouse.Name = warehouseResponseData.Data.Warehouse.Settings.Name
	}

	return &warehouseResponseData.Data.Warehouse, nil
}

func (c *Client) CreateWarehouse(enabled bool, name string, warehouseSlug string, settings WarehouseSettings) (*Warehouse, error) {
	warehouseMetadata, _ := c.GetWarehouseMetadataFromCatalog(warehouseSlug)

	newWarehouse := WarehouseRequest{
		Enabled:    enabled,
		Name:       name,
		MetadataID: &warehouseMetadata.ID,
		Settings:   settings,
	}

	newWarehouseData, err := json.Marshal(newWarehouse)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/warehouses/", c.HostURL), strings.NewReader(string(newWarehouseData)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	warehouseResponseData := WarehouseResponseData{}
	err = json.Unmarshal(body, &warehouseResponseData)
	if err != nil {
		return nil, err
	}

	return &warehouseResponseData.Data.Warehouse, nil
}

func (c *Client) UpdateWarehouse(warehouseID string, enabled bool, name string, settings WarehouseSettings) (*Warehouse, error) {
	updatedWarehouse := WarehouseRequest{
		ID:       &warehouseID,
		Enabled:  enabled,
		Name:     name,
		Settings: settings,
	}

	updatedWarehouseData, err := json.Marshal(updatedWarehouse)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/warehouses/%s", c.HostURL, warehouseID), strings.NewReader(string(updatedWarehouseData)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	warehouseResponseData := WarehouseResponseData{}
	err = json.Unmarshal(body, &warehouseResponseData)
	if err != nil {
		return nil, err
	}

	return &warehouseResponseData.Data.Warehouse, nil
}

func (c *Client) DeleteWarehouse(warehouseID string) (string, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/warehouses/%s", c.HostURL, warehouseID), nil)
	if err != nil {
		return "", err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return "", err
	}

	return "", err
}
