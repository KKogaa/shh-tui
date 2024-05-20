package ui

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/KKogaa/shh-tui/client"
	"github.com/KKogaa/shh-tui/config"
	"github.com/KKogaa/shh-tui/entity"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	viewport   viewport.Model
	messages   []string
	textarea   textarea.Model
	err        error
	wsClient   client.WSClient
	restClient client.RestClient
	config     *config.Config
}

func splitEveryNChars(str string, n int) ([]string, error) {
	if n <= 0 {
		return nil, errors.New("error wrong split for n must not be zero")
	}
	substrings := make([]string, 0, (len(str)+n-1)/n)
	for i := 0; i < len(str); i += n {
		end := min(i+n, len(str))
		substrings = append(substrings, str[i:end])
	}
	return substrings, nil
}

func CreateDisplayMsg(username string, msg string, width int) string {
	entireMsg := fmt.Sprintf("%s: %s", username, msg)
	splitSentences, err := splitEveryNChars(entireMsg, width-10)
	if err != nil {
		return entireMsg
	}
	newMsg := strings.Join(splitSentences, "\n")
	newMsg = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Render(newMsg)
	return newMsg
}

func InitialModel(restClient client.RestClient, wsClient client.WSClient,
	config *config.Config) model {

	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	ta.Prompt = "| "
	ta.CharLimit = config.Client.Chatbox.Width * config.Client.Chatbox.Height

	ta.SetWidth(config.Client.Chatbox.Width)
	ta.SetHeight(config.Client.Chatbox.Height)
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(config.Client.Chat.Width, config.Client.Chat.Height)
	vp.KeyMap.Up.SetKeys("up")
	vp.KeyMap.Down.SetKeys("down")

	respMsgs, err := restClient.GetMessagesByChatroomName(config.Client.Chatroom)
	if err != nil {
		log.Fatal(err)
	}

	messages := []string{}
	for _, msg := range respMsgs {
		displayMsg := CreateDisplayMsg(msg.User, msg.Msg, config.Client.Chat.Width)
		messages = append(messages, displayMsg)
		vp.SetContent(strings.Join(messages, "\n"))
	}

	vp.GotoBottom()

	return model{
		textarea:   ta,
		messages:   messages,
		viewport:   vp,
		err:        nil,
		wsClient:   wsClient,
		restClient: restClient,
		config:     config,
	}

}

func (m model) Init() tea.Cmd {
	return nil
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
			return m, tea.Quit
		case tea.KeyEnter:
			m.wsClient.SendMessage(m.textarea.Value())
			return m, nil
		default:
			return m, nil
		}

	case entity.Message:
		displayMsg := CreateDisplayMsg(msg.User, msg.Msg, m.config.Client.Chat.Width)
		m.messages = append(m.messages, displayMsg)
		m.viewport.SetContent(strings.Join(m.messages, "\n"))

		if msg.User == m.config.Client.Username {
			m.textarea.Reset()
		}

		m.viewport.GotoBottom()
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf("%s\n\n%s\n", m.viewport.View(), m.textarea.View())
}
