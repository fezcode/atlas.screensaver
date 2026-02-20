package waves

import (
	"math"
	"strings"
	"time"

	"atlas.screensaver/internal/savers"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	width, height int
	tickCmd       tea.Cmd
	time          float64
}

func NewModel() savers.Saver {
	return &Model{
		tickCmd: tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
			return tickMsg(t)
		}),
	}
}

type tickMsg time.Time

func (m *Model) Init() tea.Cmd {
	return m.tickCmd
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		m.time += 0.1
		return m, tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}
	return m, nil
}

func (m *Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	chars := []string{" ", ".", "░", "▒", "▓", "█"}
	var sb strings.Builder

	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			fx := float64(x) / float64(m.width)
			fy := float64(y) / float64(m.height)
			
			val := math.Sin(fx*10.0+m.time) + math.Cos(fy*10.0+m.time) + math.Sin((fx+fy)*5.0+m.time*0.5)
			val = (val + 3.0) / 6.0 // Normalise to 0-1
			
			idx := int(val * float64(len(chars)-1))
			if idx < 0 { idx = 0 }
			if idx >= len(chars) { idx = len(chars) - 1 }
			
			color := lipgloss.Color("#555555")
			if idx > 4 {
				color = lipgloss.Color("#FFD700")
			} else if idx > 2 {
				color = lipgloss.Color("#CCCCCC")
			}
			
			sb.WriteString(lipgloss.NewStyle().Foreground(color).Render(chars[idx]))
		}
		if y < m.height-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
