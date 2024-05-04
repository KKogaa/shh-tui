package main

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

func main() {

	// conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	conn, _, err := websocket.DefaultDialer.Dial("ws://ec2-3-84-186-150.compute-1.amazonaws.com:8080/ws", nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}

	p := tea.NewProgram(initialModel(conn))

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("error reading from websocket")
				return
			}

			var decodedMsg Message
			err = json.Unmarshal(message, &decodedMsg)
			if err != nil {
				log.Println("error decoding json")
				return
			}

			p.Send(decodedMsg)
		}
	}()

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	viewport      viewport.Model
	messages      []string
	textarea      textarea.Model
	senderStyle   lipgloss.Style
	err           error
	websocketConn *websocket.Conn
}

func initialModel(conn *websocket.Conn) model {

	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	ta.Prompt = "| "
	ta.CharLimit = 280

	//TODO: resize this model depending on the terminal window
	ta.SetWidth(100)
	ta.SetHeight(3)
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(100, 5)
	// vp.SetContent("")

	return model{
		textarea:      ta,
		messages:      []string{},
		viewport:      vp,
		senderStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:           nil,
		websocketConn: conn,
	}

}

func (m model) Init() tea.Cmd {
	return nil
}

// TODO: mimick getting the data each time from a fake producer
type Message struct {
	Msg  string `json:"msg"`
	User string `json:"user"`
	// Timestamp time.Time `json:"timestamp"`
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			msg := Message{
				User: "test",
				Msg:  m.textarea.Value(),
			}
			jsonMsg, err := json.Marshal(msg)
			err = m.websocketConn.WriteMessage(websocket.TextMessage, []byte(jsonMsg))
			if err != nil {
				log.Println("Error sending message to WebSocket server:", err)
				//TODO: add something that tells that message error to send
				return m, nil
			}
			// m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
			// m.viewport.SetContent(strings.Join(m.messages, "\n"))
			// m.textarea.Reset()
			// m.viewport.GotoBottom()
			return m, nil
		}
	case Message:
		displayMsg := fmt.Sprintf("%s: %s", msg.User, msg.Msg)
		m.messages = append(m.messages, m.senderStyle.Render(displayMsg))
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.textarea.Reset()
		m.viewport.GotoBottom()
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf("%s\n%s\n", m.viewport.View(), m.textarea.View())
}
