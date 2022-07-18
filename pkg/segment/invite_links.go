package segment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Permission struct {
	RoleID string `json:"roleId"`
}

type Invite struct {
	Email       string       `json:"email"`
	Permissions []Permission `json:"permissions"`
}

type CreateInviteLinksRequestData struct {
	Invites []Invite `json:"invites"`
}

type Emails struct {
	Emails []string `json:"emails"`
}

type CreateInviteLinksResponseData struct {
	Data Emails `json:"data"`
}

func (c *Client) CreateInviteLink(email string, roleID string) (string, error) {
	permission := Permission{
		RoleID: roleID,
	}
	invite := Invite{
		Email:       email,
		Permissions: []Permission{permission},
	}
	createInviteLinksRequestData := CreateInviteLinksRequestData{
		Invites: []Invite{invite},
	}

	requestData, err := json.Marshal(createInviteLinksRequestData)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/invites", c.HostURL), strings.NewReader(string(requestData)))
	if err != nil {
		return "", err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return "", err
	}

	createInviteLinksResponseData := CreateInviteLinksResponseData{}
	err = json.Unmarshal(body, &createInviteLinksResponseData)
	if err != nil {
		return "", err
	}

	return createInviteLinksResponseData.Data.Emails[0], nil
}

func (c *Client) DeleteInviteLink(email string) (string, error) {
	emails := Emails{
		Emails: []string{email},
	}

	requestData, err := json.Marshal(emails)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/invites", c.HostURL), strings.NewReader(string(requestData)))
	if err != nil {
		return "", err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return "", err
	}

	return "", err
}
