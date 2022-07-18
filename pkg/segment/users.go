package segment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Resource struct {
	ID     string  `json:"id"`
	Type   string  `json:"type"`
	Labels []Label `json:"labels"`
}

type User struct {
	ID          *string      `json:"id,omitempty"`
	Email       string       `json:"email"`
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
}

type UserResponse struct {
	User User `json:"user"`
}

type UserResponseData struct {
	Data UserResponse `json:"data"`
}

type UserIDs struct {
	UserIDs []string `json:"userIds"`
}


func (c *Client) GetUser(userID string) (*User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/%s", c.HostURL, userID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	userResponseData := UserResponseData{}
	err = json.Unmarshal(body, &userResponseData)
	if err != nil {
		return nil, err
	}

	return &userResponseData.Data.User, nil
}

func (c *Client) UpdateUser(userID string) (*User, error) {
	// Currently doesn't do anything as the API docs are a little strange
	user, err := c.GetUser(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (c *Client) DeleteUser(userID string, email string) (string, error) {
	user, err := c.GetUser(userID)
	if user != nil {
		userIDs := UserIDs{[]string{userID}}
		requestData, err := json.Marshal(userIDs)
		if err != nil {
			return "", err
		}

		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/users", c.HostURL), strings.NewReader(string(requestData)))
		if err != nil {
			return "", err
		}

		_, err = c.doRequest(req)
		if err != nil {
			return "", err
		}
		return "", err
	}

	_, err = c.DeleteInviteLink(email)

	if err != nil {
		return "", err
	}

	return "", err
}
