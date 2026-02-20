package dna

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
	offset        float64
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
		m.offset += 0.2
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

	var sb strings.Builder
	cx := float64(m.width) / 2
	amplitude := 10.0
	if cx < amplitude {
		amplitude = cx - 2
	}

	for y := 0; y < m.height; y++ {
		t := m.offset + float64(y)*0.3
		x1 := cx + math.Sin(t)*amplitude
		x2 := cx + math.Sin(t+math.Pi)*amplitude

		row := make([]rune, m.width)
		for i := range row {
			row[i] = ' '
		}

		// Connectors
		start := int(x1)
		end := int(x2)
		if start > end {
			start, end = end, start
		}
		for i := start + 1; i < end; i++ {
			if i >= 0 && i < m.width {
				if y%4 == 0 {
					row[i] = '-'
				}
			}
		}

		// Nodes
		c1 := lipgloss.Color("#FFD700")
		c2 := lipgloss.Color("#00D7FF")
		
		s1 := lipgloss.NewStyle().Foreground(c1).Bold(true).Render("8")
		s2 := lipgloss.NewStyle().Foreground(c2).Bold(true).Render("8")

		var rowStr strings.Builder
		for i, r := range row {
			if i == int(x1) {
				rowStr.WriteString(s1)
			} else if i == int(x2) {
				rowStr.WriteString(s2)
			} else {
				rowStr.WriteRune(r)
			}
		}
		sb.WriteString(rowStr.String())
		if y < m.height-1 {
			sb.WriteString("
")
		}
	}

	return sb.String()
}
