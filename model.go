package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jmorganca/ollama/api"
	"golang.org/x/term"
	"terminal-gpt/colors"
	"terminal-gpt/utils"
)

const PromptTemplate = "Which terminal command best describes the following: %s. respond with an object with 3 values \"command\" which is the suggested command. \"description\" which explains what the command does and \"safe\" which is either true or false depending on if the command could be dangerous to execute automatically."

type Model struct {
	Input          textinput.Model
	Loading        bool
	Response       OllamaResponse
	Actions        []string
	selectedAction int
	spinner        spinner.Model
	spinnerTick    int
	bannerText     string
	bannerStyle    lipgloss.Style
}

type OllamaResponse struct {
	Command     string
	Description string
	Safe        bool
}

func SendResponse(data OllamaResponse) {
	Program.Send(ResponseMsg{data})
}

type ResponseMsg struct {
	Response OllamaResponse
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

func (m Model) successBanner(content string) Model {
	m.bannerStyle = lipgloss.NewStyle().
		Background(colors.Success).
		Padding(0, 1).
		Foreground(colors.SuccessContent)
	m.bannerText = content
	return m
}

func (m Model) warningBanner(content string) Model {
	m.bannerStyle = lipgloss.NewStyle().
		Background(colors.Danger).
		Padding(0, 1).
		Foreground(colors.DangerContent)
	m.bannerText = content
	return m
}

func (m Model) resetBanner() Model {
	m.bannerText = ""
	m.bannerStyle = lipgloss.NewStyle()
	return m
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		switch msg.(tea.KeyMsg).String() {
		case "up", "k":
			m.selectedAction--
			if m.selectedAction < 0 {
				m.selectedAction = 3
			}
		case "down", "j":
			m.selectedAction++
			if m.selectedAction > 3 {
				m.selectedAction = 0
			}
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.Input.Focused() {
				m.Input.Blur()
				m.SendRequest()
				m.resetBanner()
				m.Loading = true
				return m, m.spinner.Tick
			} else {
				switch m.selectedAction {
				case 0:
					utils.CopyToClipboard(m.Response.Command)
					return m.successBanner("Copied to clipboard"), nil
				case 1:
					if m.Response.Safe {
						utils.ExecuteGeneratedCommand(m.Response.Command)
						return m.successBanner("Command executed"), nil
					} else {
						return m.warningBanner("Command not executed because it is not safe, copy it and run it manually"), nil
					}
				case 2:
					m.Input.Focus()
					m.Input.SetValue("")
					m.Response = OllamaResponse{}
					return m.resetBanner(), nil
				case 3:
					return m, tea.Quit
				}
			}
		}
	case ResponseMsg:
		m.Loading = false
		m.Input.Blur()
		m.Response = msg.Response
		if m.Response.Safe {
			return m.resetBanner(), nil
		} else {
			return m.warningBanner("WARNING: This command could be dangerous to run"), nil
		}
	case spinner.TickMsg:
		m.spinnerTick++
		if m.spinnerTick > len(m.spinner.Spinner.Frames)*2 {
			m.spinnerTick = 0
			m.spinner = spinner.New(spinner.WithSpinner(randomSpinner()))
			return m, m.spinner.Tick
		}
	case bannerMsg:
		m.bannerText = msg.Text
		return m, nil
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	var cmd2 tea.Cmd
	m.Input, cmd2 = m.Input.Update(msg)
	return m, tea.Batch(cmd, cmd2)
}

type bannerMsg struct {
	Text string
}

func SetBanner(text string) {
	Program.Send(
		bannerMsg{
			Text: text,
		},
	)
}

func (m Model) View() string {
	w, h, _ := term.GetSize(0)
	if m.Loading == true {
		return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, m.viewLoading())
	} else if m.Input.Focused() {
		app := lipgloss.NewStyle().
			Background(colors.Primary).
			Foreground(colors.PrimaryContent).
			Padding(1, 1, 1, 1).
			Border(lipgloss.InnerHalfBlockBorder()).
			BorderForeground(colors.Primary).
			Render(lipgloss.JoinVertical(lipgloss.Left, "Input:", m.Input.View()))
		appWidth := lipgloss.Width(app)
		if m.bannerText != "" {
			app = lipgloss.JoinVertical(lipgloss.Left, m.bannerStyle.Width(appWidth).Render(m.bannerText), app)
		}
		return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, app)
	} else {
		return m.showResponse()
	}
}

func (m Model) viewLoading() string {
	return lipgloss.JoinVertical(lipgloss.Center, m.spinner.View(), m.bannerStyle.Render(m.bannerText))
}

func (m Model) SendRequest() {
	client, e := api.ClientFromEnvironment()
	if e != nil {
		panic(e)
	}
	val := m.Input.Value()
	m.Input.SetValue("")
	m.Input.Blur()

	go func() {
		stream := false
		ctx := context.Background()
		_, err := client.Show(ctx, &api.ShowRequest{
			Model: "gemma:2b",
		})
		if err != nil {
			SetBanner("Pulling the model from the server")
			client.Pull(ctx, &api.PullRequest{
				Name: "gemma:2b",
			}, func(resp api.ProgressResponse) error {
				total := resp.Total
				completed := resp.Completed
				percent := int((float64(completed) / float64(total)) * 100)
				SetBanner(fmt.Sprintf("Pulling model: %d%%", percent))
				return nil
			})
		}
		e = client.Generate(ctx, &api.GenerateRequest{
			Model:  "gemma:2b",
			Prompt: fmt.Sprintf(PromptTemplate, val),
			Format: "json",
			Stream: &stream,
		}, func(resp api.GenerateResponse) error {
			var r OllamaResponse
			e := json.Unmarshal([]byte(resp.Response), &r)
			if e != nil {
				panic(e)
			}
			SendResponse(r)
			return nil
		})
		if e != nil {
			panic(e)
		}
	}()
}

func (m Model) showResponse() string {
	w, h, _ := term.GetSize(0)

	responseWidth := w / 2
	responseMaxWidth := 40
	if responseWidth > responseMaxWidth {
		responseWidth = responseMaxWidth

	}
	actionsWidth := w / 3
	actionsMaxWidth := 20
	if actionsWidth > actionsMaxWidth {
		actionsWidth = actionsMaxWidth
	}

	labelStyle := s().Bold(false).Italic(true).Background(colors.Primary).Faint(true)
	var responseParts [][]any
	responseParts = append(responseParts, []any{labelStyle, "Suggested command:"})
	responseParts = append(responseParts, []any{s(), m.Response.Command})
	responseParts = append(responseParts, []any{labelStyle, "Description:"})
	responseParts = append(responseParts, []any{s(), m.Response.Description})

	responsePartStrings := make([]string, len(responseParts))
	for i, part := range responseParts {
		responsePartStrings[i] = part[1].(string)
	}
	responseWidth = getWidestString(responsePartStrings)
	if responseWidth > responseMaxWidth {
		responseWidth = responseMaxWidth
	}

	responsePartsStrings := make([]string, len(responseParts))
	for i, part := range responseParts {
		responsePartsStrings[i] = part[0].(lipgloss.Style).Width(responseWidth).Render(part[1].(string))
	}

	var actionMenu []string
	for i, action := range m.Actions {
		st := s().Width(actionsWidth)
		if i == m.selectedAction {
			st = st.Background(colors.Secondary).Foreground(colors.SecondaryContent)
		}
		actionMenu = append(actionMenu, st.Render(action))
	}

	responsePartsHeight := lipgloss.Height(lipgloss.JoinVertical(lipgloss.Left, responsePartsStrings...))
	actionMenuHeight := lipgloss.Height(lipgloss.JoinVertical(lipgloss.Left, actionMenu...))
	var ch int
	if responsePartsHeight > actionMenuHeight {
		ch = responsePartsHeight + 2
	} else {
		ch = actionMenuHeight + 2
	}

	partResponse := lipgloss.NewStyle().
		Background(colors.Primary).
		Foreground(colors.PrimaryContent).
		Padding(1, 1, 1, 1).
		Border(lipgloss.InnerHalfBlockBorder()).
		BorderForeground(colors.Primary).
		Height(ch).
		Render(lipgloss.JoinVertical(lipgloss.Left, responsePartsStrings...))
	partActions := lipgloss.NewStyle().
		Background(colors.Primary).
		Foreground(colors.PrimaryContent).
		Padding(1, 1, 1, 1).
		Border(lipgloss.InnerHalfBlockBorder()).
		BorderForeground(colors.Primary).
		Height(ch).
		Render(lipgloss.JoinVertical(lipgloss.Left, actionMenu...))

	app := lipgloss.JoinHorizontal(lipgloss.Left,
		partResponse,
		partActions,
	)
	appWidth := lipgloss.Width(app)
	if m.bannerText != "" {
		app = lipgloss.JoinVertical(lipgloss.Left, m.bannerStyle.Width(appWidth).Render(m.bannerText), app)
	}

	return lipgloss.Place(
		w,
		h,
		lipgloss.Center,
		lipgloss.Center,
		app,
	)

}

func getWidestString(strings []string) int {
	w := 0
	for _, s := range strings {
		if len(s) > w {
			w = len(s)
		}
	}
	return w
}
