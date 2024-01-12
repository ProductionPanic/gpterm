package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
	"math/rand"
	"strings"
	"terminal-gpt/colors"
)

func randomSpinner() spinner.Spinner {
	spinners := []spinner.Spinner{
		spinner.Dot,
		spinner.Ellipsis,
		spinner.Line,
		spinner.Globe,
		spinner.Moon,
		spinner.Monkey,
		spinner.Jump,
		spinner.MiniDot,
		spinner.Points,
		spinner.Pulse,
		spinner.Meter,
		spinner.Hamburger,
	}
	return spinners[rand.Intn(len(spinners))]
}

var Program *tea.Program

func main() {
	w, _, _ := term.GetSize(0)
	in := textinput.New()
	inw := 60
	if w-4 < inw {
		inw = w - 4
	}
	in.Width = inw
	placeholder := "Describe the command you want to run"
	if len(placeholder) < inw {
		dif := inw - len(placeholder) + 1
		placeholder = placeholder + strings.Repeat(" ", dif)
	}
	in.Placeholder = placeholder
	in.Focus()
	in.Prompt = ""
	in.PromptStyle = lipgloss.NewStyle().Foreground(colors.PrimaryContent)
	in.TextStyle = lipgloss.NewStyle().Foreground(colors.Primary).Background(colors.PrimaryContent)
	in.PlaceholderStyle = in.TextStyle
	in.Cursor.Style = lipgloss.NewStyle().Foreground(colors.Secondary).Background(colors.SecondaryContent)
	s := spinner.New(spinner.WithSpinner(randomSpinner()))
	model := Model{
		Loading:        false,
		Input:          in,
		Actions:        []string{"Copy to clipboard", "Run command", "Ask again", "Quit"},
		spinner:        s,
		selectedAction: 0,
		bannerText:     "",
		bannerStyle:    lipgloss.NewStyle(),
	}
	Program = tea.NewProgram(model)
	_, e := Program.Run()
	if e != nil {
		panic(e)
	}
}

func s() lipgloss.Style {
	return lipgloss.NewStyle()
}
