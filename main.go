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

/*
*
styling
*
*/
var (
	columnStyle  = lipgloss.NewStyle().Padding(1, 2)
	focusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("62")).
			Border(lipgloss.RoundedBorder())
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

type Task struct {
	status      status
	title       string
	description string
}

func (t *Task) Next() {
	if t.status == done {
		t.status = todo
	} else {
		t.status++
	}
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

func (m *Model) Next() {
	if m.focused == done {
		m.focused = todo
	} else {
		m.focused++
	}
}

func (m *Model) Prev() {
	if m.focused == todo {
		m.focused = done
	} else {
		m.focused--
	}
}

func (m *Model) MoveToNext() tea.Msg {
	selectedItem := m.lists[m.focused].SelectedItem()
	selectedTask := selectedItem.(Task)
	m.lists[selectedTask.status].RemoveItem(m.lists[m.focused].Index())
	selectedTask.Next()
	m.lists[selectedTask.status].InsertItem(len(m.lists[selectedTask.status].Items())-1, list.Item(selectedTask))
	return nil
}

func (m *Model) DeleteTask() tea.Msg {
	m.lists[m.focused].RemoveItem(m.lists[m.focused].Index())
	return nil
}

func (m *Model) initList(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height-divisor/2)
	defaultList.SetShowHelp(false)

	m.lists = []list.Model{defaultList, defaultList, defaultList}
	//init todos
	m.lists[todo].Title = "Options"
	m.lists[todo].SetItems([]list.Item{
		Task{status: todo, title: "Invade Taiwan", description: "You show great spirit. Nom Nom the chips"},
		Task{status: inProgress, title: "Invade Russia", description: "Borscht and conquest"},
		Task{status: done, title: "Invade Japan", description: "All their islands are belong to us."},
	})
	//in progress
	m.lists[inProgress].Title = "In Progress"
	m.lists[inProgress].SetItems([]list.Item{
		Task{status: todo, title: "Drink Tea", description: "Good leaves."},
		Task{status: inProgress, title: "Make plans", description: "One year, two year, ten year, thirty year."},
		Task{status: done, title: "Make Boats", description: "Boats are good."},
	})
	//done
	m.lists[done].Title = "Done"
	m.lists[done].SetItems([]list.Item{
		Task{status: todo, title: "Eat Food", description: "Nom Nom the chips"},
		Task{status: inProgress, title: "Build Wall", description: "Because."},
		Task{status: done, title: "Have flood", description: "So much water."},
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
		case "left", "h", "shift+tab":
			m.Prev()
		case "right", "l", "tab":
			m.Next()
		case "enter":
			m.MoveToNext()
		case "backspace":
			m.DeleteTask()
		}

	case tea.WindowSizeMsg:
		if !m.loaded {
			m.initList(msg.Width, msg.Height)
			m.loaded = true
		}
		m.initList(msg.Width, msg.Height)
	}
	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.loaded {
		todoView := m.lists[todo].View()
		inProgressView := m.lists[inProgress].View()
		doneView := m.lists[done].View()
		//switch
		switch m.focused {
		case inProgress:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				columnStyle.Render(todoView),
				focusedStyle.Render(inProgressView),
				columnStyle.Render(doneView))
		case done:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				columnStyle.Render(todoView),
				columnStyle.Render(inProgressView),
				focusedStyle.Render(doneView))

		default:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				focusedStyle.Render(todoView),
				columnStyle.Render(inProgressView),
				columnStyle.Render(doneView))
		}
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
