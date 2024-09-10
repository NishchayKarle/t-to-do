# T-To-Do - Terminal To-do app

T-To-Do is a simple, terminal-based task manager built using [Bubble Tea](https://github.com/charmbracelet/bubbletea). It allows you to manage a list of tasks with keyboard shortcuts for quick and efficient task management. You can save tasks to a file, update them, add subtasks, and mark tasks as completed—all without leaving your terminal.

![example.gif](/T-To-Do-Example.gif)

## Features
- Create and delete tasks
- Indent tasks to create subtasks
- Mark tasks as completed or undone
- Save tasks to a file automatically
- Key bindings for fast navigation
- Helpful display of available commands
- Quit and save or quit without saving

## Key Bindings

| Key             | Action                                |
|-----------------|---------------------------------------|
| `h`             | Show/hide help                        |
| `n`             | Add a new task                        |
| `k` / `j`       | Move cursor up/down                   |
| `tab` / `shift+tab` | Increase/decrease task indentation |
| `enter` / `space` | Mark task as done/undone            |
| `alt+enter`     | Save task and add new task            |
| `d`             | Delete task                           |
| `g`             | Move to top of the task list          |
| `G`             | Move to bottom of the task list       |
| `{` / `}`       | Move up/down 3 tasks                  |
| `ctrl+c` / `q`  | Save and quit                         |
| `Q`             | Quit without saving                   |