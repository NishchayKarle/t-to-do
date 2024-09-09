package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"main.go/task"
)

const tab string = "    "

type model struct {
	fileName          string           // File name to save tasks
	CurrentTask       task.Task        // The current task
	Tasks             []string         // List of tasks
	cursor            int              // The cursor is the index of the task that is selected
	Completed         map[int]struct{} // Set of tasks that are done
	IndentationLevels []int            // Indentation level
}

func NewModel() model {
	return model{
		fileName:          "",
		CurrentTask:       task.NewTask(),
		Tasks:             make([]string, 0),
		Completed:         make(map[int]struct{}),
		cursor:            0,
		IndentationLevels: []int{0},
	}
}

func (m model) Init() tea.Cmd {
	return load
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	handleTab := func() {
		m.IndentationLevels[m.cursor]++
	}

	handleShiftTab := func() {
		if m.IndentationLevels[m.cursor] > 0 {
			m.IndentationLevels[m.cursor]--
		}
	}

	var cmd tea.Cmd

	if m.CurrentTask.NewTaskInput {
		m.CurrentTask.TextInput, cmd = m.CurrentTask.TextInput.Update(msg)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, m.saveAndQuit()

			case "tab":
				handleTab()

			case "shift+tab":
				handleShiftTab()

			case "enter":
				m.CurrentTask.NewTaskInput = false
				m.Tasks = append(m.Tasks, m.CurrentTask.TextInput.Value())
			}
		}

		return m, cmd
	}

	switch msg := msg.(type) {

	case fileMsg:
		m.fileName = string(msg)

		// Load tasks from file
		data, _ := os.ReadFile(m.fileName)
		json.Unmarshal(data, &m)

		if len(m.Tasks) == 0 {
			m.CurrentTask.NewTaskInput = true
			return m, textinput.Blink
		}

	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, m.saveAndQuit()

		case "n":
			m.CurrentTask.NewTaskInput = true
			m.CurrentTask.TextInput.Reset()
			m.IndentationLevels = append(m.IndentationLevels, 0)
			m.cursor = len(m.Tasks)

		case "j", "down":
			if m.cursor < len(m.Tasks)-1 {
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
			_, ok := m.Completed[m.cursor]
			if ok {
				delete(m.Completed, m.cursor)
			} else {
				m.Completed[m.cursor] = struct{}{}
			}

		}

	}
	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("Tasks: (%s):\n\n", m.fileName)

	for i, currentTask := range m.Tasks {
		// Set cursor position
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		// Check if task is completed
		checked := " "
		if _, ok := m.Completed[i]; ok {
			checked = "x"
		}

		// Build task line with indentation
		indentation := strings.Repeat(tab, m.IndentationLevels[i])
		s += fmt.Sprintf("%s%s [%s] %s\n", indentation, cursor, checked, currentTask)
	}

	// Handle new task input if the current task is not done
	if m.CurrentTask.NewTaskInput {
		indentation := strings.Repeat(tab, m.IndentationLevels[len(m.IndentationLevels)-1])
		s += fmt.Sprintf("%s%s", indentation, m.CurrentTask.TextInput.View())
	}

	// Remove trailing newline and add footer instructions
	s = strings.TrimSuffix(s, "\n") + "\n\nPress n to add a new task.\nPress ctrl+c to quit.\n"

	return s
}

func (m model) saveAndQuit() tea.Cmd {
	if len(m.Tasks) == 0 {
		return tea.Quit
	}

	val, err := json.Marshal(m)
	if err != nil {
		fmt.Println("Error marshalling tasks:", err)
		return tea.Quit
	}

	// Save tasks to file - create file if it doesn't exist
	err = os.WriteFile(m.fileName, val, 0644)
	if err != nil {
		fmt.Println("Error saving tasks:", err)
	}
	fmt.Println("\n\nSaved tasks to", m.fileName)

	return tea.Quit
}

func load() tea.Msg {
	if len(os.Args) < 2 {
		today := time.Now().Format("01-02-2006")
		return fileMsg(today)
	} else {
		file := os.Args[1]
		return fileMsg(file)
	}
}

type fileMsg string

func main() {
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting program:", err)
		os.Exit(1)
	}
}
