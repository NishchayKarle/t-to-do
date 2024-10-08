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

const instructions = `Key Bindings:

'h'            		-> Show/hide help
'n'            		-> Add new task
'k' / 'j'      		-> Move cursor up/down
'tab' / 'shift+tab' 	-> Increase/decrease indentation
'enter' / 'space'   	-> Mark task as done
'alt+enter'        	-> Save task and add new task
'd'            		-> Delete task
'g'            		-> Move to top
'G'            		-> Move to bottom
'{' / '}'      		-> Move up/down 3 tasks
'ctrl+c' / 'q' 		-> Quit
`

type model struct {
	fileName          string           // File name to save tasks
	CurrentTask       task.Task        // The current task
	Tasks             []string         // List of tasks
	cursor            int              // The cursor is the index of the task that is selected
	Completed         map[int]struct{} // Set of tasks that are done
	IndentationLevels []int            // Indentation level
	showHelp          bool             // Show help
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

	cancelNewTaskInput := func() {
		m.CurrentTask.NewTaskInput = false
		m.cursor = len(m.Tasks) - 1
		m.IndentationLevels = m.IndentationLevels[:len(m.IndentationLevels)-1]
	}

	createNewTaskInput := func() {
		m.CurrentTask.NewTaskInput = true
		m.CurrentTask.TextInput.Reset()
		m.IndentationLevels = append(m.IndentationLevels, 0)
		m.cursor = len(m.Tasks)
	}

	var cmd tea.Cmd

	if m.CurrentTask.NewTaskInput {
		m.CurrentTask.TextInput, cmd = m.CurrentTask.TextInput.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				cancelNewTaskInput()
				return m, m.saveAndQuit()

			case "tab":
				handleTab()

			case "shift+tab":
				handleShiftTab()

			case "alt+enter":
				m.CurrentTask.NewTaskInput = false
				m.Tasks = append(m.Tasks, m.CurrentTask.TextInput.Value())
				createNewTaskInput()

			case "enter":
				m.CurrentTask.NewTaskInput = false
				m.Tasks = append(m.Tasks, m.CurrentTask.TextInput.Value())

			case "esc":
				cancelNewTaskInput()
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

		case "Q":
			return m, tea.Quit

		case "n":
			createNewTaskInput()

		case "d":
			if len(m.Tasks) > 0 {
				m.Tasks = append(m.Tasks[:m.cursor], m.Tasks[m.cursor+1:]...)
				delete(m.Completed, m.cursor)

				// adjust cursor values for completed tasks
				for k := range m.Completed {
					if k > m.cursor {
						delete(m.Completed, k)
						m.Completed[k-1] = struct{}{}
					}
				}

				if m.cursor == len(m.Tasks) {
					m.cursor--
				}
				m.IndentationLevels = append(m.IndentationLevels[:m.cursor], m.IndentationLevels[m.cursor+1:]...)
			}

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

		case "g":
			m.cursor = 0

		case "G":
			m.cursor = len(m.Tasks) - 1

		case "{":
			m.cursor -= 3
			if m.cursor < 0 {
				m.cursor = 0
			}

		case "}":
			m.cursor += 3
			if m.cursor >= len(m.Tasks) {
				m.cursor = len(m.Tasks) - 1
			}

		case "h":
			m.showHelp = !m.showHelp
		}

	}
	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("Tasks: %s\n\n", strings.TrimSuffix(m.fileName, ".ttdo"))

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
		indentation := strings.Repeat(tab, m.IndentationLevels[m.cursor])
		s += fmt.Sprintf("%s%s [ ] %s", indentation, ">", m.CurrentTask.TextInput.View())
	}

	// Remove trailing newline and add footer instructions
	s = strings.TrimSuffix(s, "\n")
	if m.CurrentTask.NewTaskInput {
		s += "\n\nPress 'enter' to save task. Press 'esc' to cancel\n"
	} else if !m.showHelp && !m.CurrentTask.NewTaskInput {
		s += "\n\nPress 'q' to save & quit.\nPress 'Q' to quit without saving.\nPress 'h' for help.\n"
	} else {
		s += "\n\n" + instructions
	}

	return s
}

func (m model) saveAndQuit() tea.Cmd {
	if len(m.Tasks) == 0 {
		return tea.Quit
	}

	m.CurrentTask.NewTaskInput = false
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

	return tea.Quit
}

func load() tea.Msg {
	var file string
	if len(os.Args) < 2 {
		today := time.Now().Format("01-02-2006")
		file = today + ".ttdo"
	} else {
		file = os.Args[1]
		file = strings.TrimSuffix(file, ".ttdo")
		file += ".ttdo"
	}
	return fileMsg(file)
}

type fileMsg string

func main() {
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting program:", err)
		os.Exit(1)
	}
}
