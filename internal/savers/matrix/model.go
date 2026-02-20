package matrix

import (
	"math/rand"
	"strings"
	"time"

	"atlas.screensaver/internal/savers"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Droplet struct {
	y     float64
	speed float64
	chars []rune
}

type Model struct {
	width, height int
	droplets      []*Droplet
	tickCmd       tea.Cmd
}

func NewModel() savers.Saver {
	return &Model{
		tickCmd: tea.Tick(time.Millisecond*40, func(t time.Time) tea.Msg {
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
		m.initDroplets()
		return m, nil

	case tickMsg:
		m.step()
		return m, tea.Tick(time.Millisecond*40, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}
	return m, nil
}

var matrixChars = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (m *Model) initDroplets() {
	m.droplets = make([]*Droplet, m.width)
	for i := range m.droplets {
		m.droplets[i] = m.newDroplet()
	}
}

func (m *Model) newDroplet() *Droplet {
	length := rand.Intn(15) + 5
	chars := make([]rune, length)
	for i := range chars {
		chars[i] = matrixChars[rand.Intn(len(matrixChars))]
	}
	return &Droplet{
		y:     float64(rand.Intn(m.height * 2) - m.height),
		speed: rand.Float64()*0.5 + 0.2,
		chars: chars,
	}
}

func (m *Model) step() {
	for i, d := range m.droplets {
		d.y += d.speed
		if int(d.y)-len(d.chars) > m.height {
			m.droplets[i] = m.newDroplet()
			m.droplets[i].y = -float64(len(m.droplets[i].chars))
		}
		if rand.Float64() < 0.1 {
			d.chars[rand.Intn(len(d.chars))] = matrixChars[rand.Intn(len(matrixChars))]
		}
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

	for x, d := range m.droplets {
		for i, char := range d.chars {
			y := int(d.y) - i
			if y >= 0 && y < m.height {
				style := lipgloss.NewStyle()
				if i == 0 {
					style = style.Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
				} else {
					alpha := 1.0 - float64(i)/float64(len(d.chars))
					if alpha < 0.2 {
						style = style.Foreground(lipgloss.Color("#003300"))
					} else if alpha < 0.5 {
						style = style.Foreground(lipgloss.Color("#008800"))
					} else {
						style = style.Foreground(lipgloss.Color("#00FF00"))
					}
				}
				canvas[y][x] = style.Render(string(char))
			}
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
