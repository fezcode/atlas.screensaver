package ui

import (
	"math/rand"

	"atlas.screensaver/internal/savers"
	"atlas.screensaver/internal/savers/pipes"
	"atlas.screensaver/internal/savers/stars"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700")).MarginBottom(1)
	itemStyle  = lipgloss.NewStyle().PaddingLeft(2)
	selectedStyle = lipgloss.NewStyle().PaddingLeft(0).Foreground(lipgloss.Color("#FFD700")).SetString("> ")
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).MarginTop(1)
)

type Model struct {
	choices []string
	cursor  int
	active  savers.Saver
	width   int
	height  int
	quitting bool
}

func NewModel() Model {
	return Model{
		choices: []string{"Pipes", "Stars", "Random"},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.active != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "esc" || msg.String() == "q" {
				m.active = nil
				return m, nil
			}
		case tea.WindowSizeMsg:
			m.width, m.height = msg.Width, msg.Height
		}
		
		newModel, cmd := m.active.Update(msg)
		m.active = newModel.(savers.Saver)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			choice := m.choices[m.cursor]
			if choice == "Random" {
				r := rand.Intn(2)
				if r == 0 { choice = "Pipes" } else { choice = "Stars" }
			}
			
			if choice == "Pipes" {
				m.active = pipes.NewModel()
			} else {
				m.active = stars.NewModel()
			}
			
			// Init the saver with current size
			m.active.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
			return m, m.active.Init()
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.quitting { return "" }
	if m.active != nil {
		return m.active.View()
	}

	s := titleStyle.Render("ATLAS SCREENSAVER") + "\n\n"

	for i, choice := range m.choices {
		if m.cursor == i {
			s += selectedStyle.Render(choice) + "\n"
		} else {
			s += itemStyle.Render(choice) + "\n"
		}
	}

	s += helpStyle.Render("Select a screensaver • q to quit")
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, s)
}
