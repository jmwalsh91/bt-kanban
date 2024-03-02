package main

import "github.com/charmbracelet/bubbles/list"

type status int

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
	list list.Model
	err  error
}

func (m *Model) initList() {
	m.list = list.New([]list.item{}, list.NewDefaultDelegate())
	m.list.Title = "Options"
	m.list.SetItems([]list.item{
		Task{status: todo, title: "Invade Taiwan", description: "You show great spirit. Nom Nom the chips"},
		Task{status: inProgress, title: "Invade Russia", description: "Borscht and conquest"},
		Task{status: done, title: "Invade Japan", description: "All their neko are belong to us."},
	})

}
