package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const tab string = "    "

type task struct {
	textInput textinput.Model
	done      bool
}

type model struct {
	tasks     []task           // List of tasks
	cursor    int              // The cursor is the index of the task that is selected
	completed map[int]struct{} // Set of tasks that are done
	levels    []int            // Indentation level
	err       error            // Error message
}

func newTextInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Task Name"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50
	return ti
}

func NewModel() model {
	ti := newTextInput()

	return model{
		tasks:     []task{{textInput: ti}},
		completed: make(map[int]struct{}),
		levels:    make([]int, 5),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	handleTab := func() {
		m.levels[m.cursor]++
	}

	handleShiftTab := func() {
		if m.levels[m.cursor] > 0 {
			m.levels[m.cursor]--
		}
	}

	var cmd tea.Cmd

	if !m.tasks[m.cursor].done {
		m.tasks[m.cursor].textInput, cmd = m.tasks[m.cursor].textInput.Update(msg)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit

			case "tab":
				handleTab()

			case "shift+tab":
				handleShiftTab()

			case "enter":
				m.tasks[m.cursor].done = true
			}
		}

		return m, cmd
	}

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "n":
			ti := newTextInput()
			m.tasks = append(m.tasks, task{textInput: ti, done: false})
			m.cursor = len(m.tasks) - 1

		case "j", "down":
			if m.cursor < len(m.tasks)-1 {
				m.cursor++
			}

		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "tab":
			handleTab()

		case "shift+tab":
			handleShiftTab()

		case "enter", " ":
			_, ok := m.completed[m.cursor]
			if ok {
				delete(m.completed, m.cursor)
			} else {
				m.completed[m.cursor] = struct{}{}
			}

		}

	}
	return m, nil
}

func (m model) View() string {
	date := time.Now().Format("01-02-2006")
	s := fmt.Sprintf("\nTasks for today (%s):\n\n", date)

	for i, currentTask := range m.tasks {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.completed[i]; ok {
			checked = "x"
		}

		taskLine := ""
		for j := 0; j < m.levels[i]; j++ {
			taskLine += tab
		}

		taskString := ""
		if currentTask.done {
			taskString = currentTask.textInput.Value()
			taskLine += fmt.Sprintf("%s [%s] %s", cursor, checked, taskString)
		} else {
			taskString = currentTask.textInput.View()
			taskLine += fmt.Sprintf("%s", taskString)
		}

		s += taskLine + "\n"
	}

	s += "\n\nPress n to add a new task."
	s += "\nPress ctrl+c to quit.\n"

	return s
}

func main() {
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting program:", err)
		os.Exit(1)
	}
}
