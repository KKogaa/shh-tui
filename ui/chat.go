package ui

import (
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
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
	wsClient    client.WSClient
	restClient  client.RestClient
	config      *config.Config
}

func CreateMsg(username string, msg string) string {
	//if the message is bigger than the config width go
	// displayMsg := fmt.Sprintf("%s: %s", username, msg)
	//TODO: extract message here

	return ""
}

func InitialModel(restClient client.RestClient, wsClient client.WSClient, config *config.Config) model {

	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	ta.Prompt = "| "
	ta.CharLimit = config.Client.Chatbox.Width

	//TODO: resize this model depending on the terminal window
	ta.SetWidth(config.Client.Chatbox.Width)
	ta.SetHeight(config.Client.Chatbox.Height)
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(config.Client.Chat.Width, config.Client.Chat.Height)
	vp.KeyMap.Up.SetKeys("up")
	vp.KeyMap.Down.SetKeys("down")
	// vp.KeyMap.PageUp.SetEnabled(false)
	// vp.SetContent("")

	//TODO: clean this section to only return messages
	respMsgs, err := restClient.GetMessagesByChatroomName(config.Client.Chatroom)
	if err != nil {
		log.Fatal(err)
	}

	messages := []string{}
	for _, msg := range respMsgs {
		displayMsg := fmt.Sprintf("%s: %s", msg.User, msg.Msg)
		messages = append(messages, lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Render(displayMsg))
		vp.SetContent(strings.Join(messages, "\n"))
	}

	vp.GotoBottom()

	return model{
		textarea:    ta,
		messages:    messages,
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
		wsClient:    wsClient,
		restClient:  restClient,
		config:      config,
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
			// m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
			// m.viewport.SetContent(strings.Join(m.messages, "\n"))
			// m.textarea.Reset()
			// m.viewport.GotoBottom()
			return m, nil
		default:
			return m, nil
		}

	case entity.Message:
		displayMsg := fmt.Sprintf("%s: %s", msg.User, msg.Msg)
		m.messages = append(m.messages, m.senderStyle.Render(displayMsg))
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
