package bouncing

import (
	"math/rand"
	"strings"
	"time"

	"atlas.screensaver/internal/savers"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	width, height int
	x, y          float64
	dx, dy        float64
	text          string
	color         lipgloss.Color
	tickCmd       tea.Cmd
}

func NewModel() savers.Saver {
	return &Model{
		text:  "ATLAS",
		dx:    0.4,
		dy:    0.2,
		color: lipgloss.Color("#FFD700"),
		tickCmd: tea.Tick(time.Millisecond*30, func(t time.Time) tea.Msg {
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
		m.x = float64(m.width) / 2
		m.y = float64(m.height) / 2
		return m, nil

	case tickMsg:
		m.step()
		return m, tea.Tick(time.Millisecond*30, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}
	return m, nil
}

func (m *Model) step() {
	m.x += m.dx
	m.y += m.dy

	textW := len(m.text)
	textH := 1

	hit := false
	if m.x <= 0 || int(m.x)+textW >= m.width {
		m.dx = -m.dx
		hit = true
	}
	if m.y <= 0 || int(m.y)+textH >= m.height {
		m.dy = -m.dy
		hit = true
	}

	if hit {
		colors := []string{"#FFD700", "#00D7FF", "#FF5F5F", "#5FFF5F", "#D75F00", "#AF00FF", "#FFFFFF"}
		m.color = lipgloss.Color(colors[rand.Intn(len(colors))])
	}
}

func (m *Model) View() string {
	if m.width == 0 || m.height == 0 { return "" }

	canvas := make([][]string, m.height)
	for i := range canvas {
		canvas[i] = make([]string, m.width)
		for j := range canvas[i] {
			canvas[i][j] = " "
		}
	}

	style := lipgloss.NewStyle().Foreground(m.color).Bold(true)
	renderedText := style.Render(m.text)

	var sb strings.Builder
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			if y == int(m.y) && x == int(m.x) {
				sb.WriteString(renderedText)
				x += len(m.text) - 1
			} else {
				sb.WriteString(" ")
			}
		}
		if y < m.height-1 {
			sb.WriteString("
")
		}
	}
	return sb.String()
}
