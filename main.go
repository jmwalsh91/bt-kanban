package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type status int

const divisor = 3
const (
	todo status = iota
	inProgress
	done
)

type Task struct {
	status      status
	title       string
	description string
}

func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

type Model struct {
	focused status
	lists   []list.Model
	err     error
	loaded  bool
}

func New() *Model {
	return &Model{}
}

func (m *Model) initList(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height)
	m.lists = []list.Model{defaultList, defaultList, defaultList}
	//init todos
	m.lists[todo].Title = "Options"
	m.lists[todo].SetItems([]list.Item{
		Task{status: todo, title: "Invade Taiwan", description: "You show great spirit. Nom Nom the chips"},
		Task{status: inProgress, title: "Invade Russia", description: "Borscht and conquest"},
		Task{status: done, title: "Invade Japan", description: "All their islands are belong to us."},
	})
	//in progress
	m.lists[todo].Title = "In Progress"
	m.lists[todo].SetItems([]list.Item{
		Task{status: todo, title: "Drink Tea", description: "Good leaves."},
		Task{status: inProgress, title: "Make plans", description: "One year, two year, ten year, thirty year."},
		Task{status: done, title: "Make Boats", description: "Boats are good."},
	})

}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.initList(msg.Width, msg.Height)
	}
	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.loaded {
		return lipgloss.JoinHorizontal(lipgloss.Left, m.lists[todo].View(), m.lists[inProgress].View(), m.lists[done].View())
	} else {
		return "Loading..."
	}
}

func main() {
	m := New()
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
