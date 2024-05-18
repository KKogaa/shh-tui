package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KKogaa/shh-tui/config"
	"github.com/KKogaa/shh-tui/entity"
)

type RestClient struct {
	config *config.Config
}

func NewRestClient(config *config.Config) RestClient {
	return RestClient{
		config: config,
	}
}

func (r RestClient) GetMessagesByChatroomName(chatroomName string) ([]entity.Message, error) {
	//TODO: get this from config and modify config to only use the server ip
	apiUrl := fmt.Sprintf("http://localhost:8080/messages/chatrooms/%s", chatroomName)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var msgResponses []MessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&msgResponses); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	var messages []entity.Message
	for _, msgResponse := range msgResponses {
		messages = append(messages, ToMessage(msgResponse))
	}
	return messages, nil
}
