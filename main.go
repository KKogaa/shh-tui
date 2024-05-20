package main

import (
	"log"

	"github.com/KKogaa/shh-tui/client"
	"github.com/KKogaa/shh-tui/config"
	"github.com/KKogaa/shh-tui/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	//TODO: add cli with cobra

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	//TODO: add multiple chatrooms to connect
	//TODO: add authentication to the ws client
	restClient := client.NewRestClient(config)

	wsClient, err := client.NewWsClient(config)
	if err != nil {
		log.Fatal(err)
	}

	//TODO: add read message history from rest api, load latest messages
	p := tea.NewProgram(ui.InitialModel(restClient, wsClient, config))

	//TODO: run goroutine somewhere else
	go func() {
		for {
			message, err := wsClient.ReadMessage()
			if err != nil {
				log.Println("error reading from websocket")
				return
			}
			p.Send(message)
		}
	}()

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
