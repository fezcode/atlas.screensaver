package pipes

import (
	"math/rand"
	"strings"
	"time"

	"atlas.screensaver/internal/savers"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Pipe struct {
	x, y  int
	dir   Direction
	color lipgloss.Color
	dead  bool
}

type Cell struct {
	char  string
	color lipgloss.Color
}

type Model struct {
	width, height int
	grid          [][]Cell
	pipes         []*Pipe
	tickCmd       tea.Cmd
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
		m.reset()
		return m, nil

	case tickMsg:
		m.step()
		return m, tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}
	return m, nil
}

func (m *Model) View() string {
	if m.width == 0 || len(m.grid) == 0 {
		return ""
	}

	var sb strings.Builder
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			if y < len(m.grid) && x < len(m.grid[y]) {
				c := m.grid[y][x]
				if c.char != "" {
					sb.WriteString(lipgloss.NewStyle().Foreground(c.color).Render(c.char))
				} else {
					sb.WriteRune(' ')
				}
			} else {
				sb.WriteRune(' ')
			}
		}
		if y < m.height-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func (m *Model) reset() {
	m.grid = make([][]Cell, m.height)
	for i := range m.grid {
		m.grid[i] = make([]Cell, m.width)
	}
	m.pipes = []*Pipe{}
	
	// Start with a few pipes
	for i := 0; i < 5; i++ {
		m.addPipe()
	}
}

func (m *Model) addPipe() {
	if m.width == 0 || m.height == 0 { return }
	x := rand.Intn(m.width)
	y := rand.Intn(m.height)
	dir := Direction(rand.Intn(4))
	
	// Random color from Atlas Palette + others
	colors := []string{"#FFD700", "#00D7FF", "#FF5F5F", "#5FFF5F", "#D75F00", "#AF00FF"}
	c := lipgloss.Color(colors[rand.Intn(len(colors))])
	
	m.pipes = append(m.pipes, &Pipe{x: x, y: y, dir: dir, color: c})
	// Mark start
	m.grid[y][x] = Cell{char: "O", color: c}
}

func (m *Model) step() {
	if len(m.pipes) == 0 {
		m.reset() // Restart if all dead
		return
	}

	active := 0
	for _, p := range m.pipes {
		if p.dead {
			continue
		}
		active++

		// Calculate new pos
		nx, ny := p.x, p.y
		switch p.dir {
		case Up:    ny--
		case Down:  ny++
		case Left:  nx--
		case Right: nx++
		}

		// Check bounds and collision
		if nx < 0 || nx >= m.width || ny < 0 || ny >= m.height || m.grid[ny][nx].char != "" {
			p.dead = true
			continue
		}

		// Decide new direction
		newDir := p.dir
		if rand.Float64() < 0.2 { // 20% chance to turn
			if p.dir == Up || p.dir == Down {
				if rand.Float64() < 0.5 { newDir = Left } else { newDir = Right }
			} else {
				if rand.Float64() < 0.5 { newDir = Up } else { newDir = Down }
			}
		}

		// Draw current cell
		char := ""
		if p.dir == newDir {
			if p.dir == Up || p.dir == Down { char = "║" } else { char = "═" }
		} else {
			switch p.dir {
			case Up:
				if newDir == Left { char = "╗" } else { char = "╔" }
			case Down:
				if newDir == Left { char = "╝" } else { char = "╚" }
			case Left:
				if newDir == Up { char = "╚" } else { char = "╔" }
			case Right:
				if newDir == Up { char = "╝" } else { char = "╗" }
			}
		}

		m.grid[p.y][p.x] = Cell{char: char, color: p.color}
		
		p.x, p.y = nx, ny
		p.dir = newDir
	}
	
	if active == 0 || rand.Float64() < 0.05 {
		if len(m.pipes) < 15 {
			m.addPipe()
		} else if active == 0 {
			m.reset()
		}
	}
}
