package stars

import (
	"math/rand"
	"strings"
	"time"

	"atlas.screensaver/internal/savers"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Star struct {
	x, y, z float64
	char    string
	color   lipgloss.Color
}

type Model struct {
	width, height int
	stars         []*Star
	tickCmd       tea.Cmd
}

func NewModel() savers.Saver {
	return &Model{
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
		m.initStars()
		return m, nil

	case tickMsg:
		m.step()
		return m, tea.Tick(time.Millisecond*30, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}
	return m, nil
}

func (m *Model) View() string {
	if m.width == 0 { return "" }
	
	canvas := make([][]string, m.height)
	for i := range canvas {
		canvas[i] = make([]string, m.width)
		for j := range canvas[i] {
			canvas[i][j] = " "
		}
	}

	cx, cy := float64(m.width)/2, float64(m.height)/2

	for _, s := range m.stars {
		factor := 20.0
		sx := (s.x / s.z) * factor + cx
		sy := (s.y / s.z) * factor * 0.5 + cy

		if sx >= 0 && sx < float64(m.width) && sy >= 0 && sy < float64(m.height) {
			style := lipgloss.NewStyle().Foreground(s.color)
			
			char := "."
			if s.z < 5 { char = "*" }
			if s.z < 2 { char = "@" }
			
			canvas[int(sy)][int(sx)] = style.Render(char)
		}
	}

	var sb strings.Builder
	for i, row := range canvas {
		for _, cell := range row {
			sb.WriteString(cell)
		}
		if i < m.height-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func (m *Model) initStars() {
	m.stars = make([]*Star, 100)
	for i := range m.stars {
		m.stars[i] = m.newStar()
	}
}

func (m *Model) newStar() *Star {
	x := (rand.Float64() - 0.5) * 100
	y := (rand.Float64() - 0.5) * 100
	z := rand.Float64() * 20 + 1
	
	colors := []string{"#FFFFFF", "#CCCCCC", "#888888", "#FFD700", "#5FFF5F"}
	c := lipgloss.Color(colors[rand.Intn(len(colors))])
	
	return &Star{x: x, y: y, z: z, color: c}
}

func (m *Model) step() {
	speed := 0.2
	for _, s := range m.stars {
		s.z -= speed
		if s.z <= 0.1 {
			ns := m.newStar()
			s.x, s.y, s.z, s.color = ns.x, ns.y, 20, ns.color
		}
	}
}
