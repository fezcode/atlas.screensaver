package savers

import tea "github.com/charmbracelet/bubbletea"

// Saver is just a standard Bubble Tea model.
// Implementations should handle their own animation loops (tick) and resizing.
type Saver interface {
	tea.Model
}
