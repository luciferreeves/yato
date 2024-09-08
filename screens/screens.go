package screens

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

type ScreenSwitcher struct {
	currentScreen tea.Model
}

type Globals struct {
	width  int
	height int
}

var globals Globals

func (s ScreenSwitcher) Init() tea.Cmd {
	return s.currentScreen.Init()
}

func (s ScreenSwitcher) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var model tea.Model

	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		globals.width, globals.height = m.Width, m.Height
	}

	model, cmd = s.currentScreen.Update(msg)

	return ScreenSwitcher{currentScreen: model}, cmd
}

func (s ScreenSwitcher) View() string {
	return s.currentScreen.View()
}

func (s ScreenSwitcher) Switch(screen tea.Model) (tea.Model, tea.Cmd) {
	s.currentScreen = screen
	return s.currentScreen, s.currentScreen.Init()
}

func screen() ScreenSwitcher {
	screen := homeScreen()

	return ScreenSwitcher{
		currentScreen: screen,
	}
}

func Initialize() tea.Model {

	width, height, err := term.GetSize(int(os.Stdout.Fd()))

	if err != nil {
		width = 80
		height = 30
	}

	globals.width = width
	globals.height = height

	return screen()
}
