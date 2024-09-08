package screens

import tea "github.com/charmbracelet/bubbletea"

type HomeScreen struct {
}

func homeScreen() tea.Model {
	return HomeScreen{}
}

func (h HomeScreen) Init() tea.Cmd {
	return nil
}

func (h HomeScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if msg.Type == tea.KeyCtrlC {
			return h, tea.Quit
		}
	}

	return h, nil
}

func (h HomeScreen) View() string {
	return "Home"
}
