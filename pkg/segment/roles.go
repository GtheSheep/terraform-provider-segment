package segment

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Role struct {
	ID          *string `json:"id,omitempty"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
}

type RolesResponse struct {
	Roles []Role `json:"roles"`
}

type RoleResponseData struct {
	Data RolesResponse `json:"data"`
}

func (c *Client) GetRole(roleName string) (*Role, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/roles?pagination.count=100", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	roleResponseData := RoleResponseData{}
	err = json.Unmarshal(body, &roleResponseData)
	if err != nil {
		return nil, err
	}

	for i, role := range roleResponseData.Data.Roles {
		if role.Name == roleName {
			return &roleResponseData.Data.Roles[i], nil
		}
	}

	return nil, fmt.Errorf("Did not find role %s", roleName)
}
