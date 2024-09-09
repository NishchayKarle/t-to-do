package task

import "github.com/charmbracelet/bubbles/textinput"

func newTextInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Task Name"
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 100
	ti.Prompt = ""
	return ti
}

type Task struct {
	TextInput    textinput.Model `json:"-"`
	NewTaskInput bool            `json:"NewTaskInput"`
}

func NewTask() Task {
	ti := newTextInput()
	return Task{
		TextInput:    ti,
		NewTaskInput: false,
	}
}
