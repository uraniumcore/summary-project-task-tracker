package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type UpdateDesc int

const (
	StatusChange UpdateDesc = iota
	TaskChange
)

type Status string

const (
	StatusTodo       = "todo"
	StatusInProgress = "in-progress"
	StatusDone       = "done"
)

type Task struct {
	ID          int            `json:"id"`
	Description string         `json:"description"`
	Status      Status         `json:"status"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	History     []UpdateRecord `json:"history"`
}

type UpdateRecord struct {
	UpdateDesc UpdateDesc `json:"update_desc"`
	UpdateTime time.Time  `json:"update_time"`
}

func (u UpdateDesc) String() string {
	switch u {
	case StatusChange:
		return "StatusChange"
	case TaskChange:
		return "TaskChange"
	default:
		return "UnknownChange"
	}
}

func addTask(task string) (string, error) {
	tasks, err := readJSON("tasks.json")
	if err != nil {
		return "", err
	}

	lastID := nextID(tasks)

	tasks = append(tasks, Task{ID: lastID + 1, Description: task, Status: StatusTodo, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		History: []UpdateRecord{{UpdateDesc: TaskChange, UpdateTime: time.Now()}}})

	err = updateJSON("tasks.json", tasks)
	if err != nil {
		return "", err
	}

	return "task successfully added", nil
}

func listTasks(arg string) (string, error) {
	var count int
	tasks, err := readJSON("tasks.json")
	if err != nil {
		return "", err
	}

	if arg == "all" || strings.TrimSpace(arg) == "" {
		for _, task := range tasks {
			fmt.Println("\nDescription ID:", task.ID, "\nDescription:", task.Description, "\nStatus:", task.Status, "\nCreated at:", task.CreatedAt, "\nUpdated at:", task.UpdatedAt, "\n")
			count++
		}
		fmt.Println("Filtered tasks:", count)
	} else {
		filter, err := ParseTaskStatus(arg)

		if err != nil {
			log.Println(err)
			return "", err
		}

		for _, task := range tasks {
			if task.Status == filter {
				fmt.Println("\nDescription ID:", task.ID, "\nDescription:", task.Description, "\nStatus:", task.Status, "\nCreated at:", task.CreatedAt, "\nUpdated at:", task.UpdatedAt, "\n")
				count++
			}
		}
		fmt.Println("Filtered tasks:", count)
	}
	msg := `Total tasks: ` + fmt.Sprint(len(tasks))
	return msg, nil
}

func updateTask(idS, newTask string) (string, error) {
	tasks, err := readJSON("tasks.json")
	if err != nil {
		return "", err
	}

	id, err := strconv.Atoi(idS)
	if err != nil {
		return "", err
	}

	found := false

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Description = newTask
			newUpdateTask(tasks, i, TaskChange, time.Now())
			found = true
			break
		}
	}

	if !found {
		return "", errors.New("task doesn't exist")
	}

	err = updateJSON("tasks.json", tasks)
	if err != nil {
		return "", err
	}

	return "task successfully updated", nil
}

func deleteTask(idS string) (string, error) {
	tasks, err := readJSON("tasks.json")
	if err != nil {
		return "", err
	}

	id, err := strconv.Atoi(idS)
	if err != nil {
		return "", err
	}

	found := false

	for i := range tasks {
		if tasks[i].ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return "", errors.New("task doesn't exist")
	}

	err = updateJSON("tasks.json", tasks)
	if err != nil {
		return "", err
	}

	return "task successfully deleted", nil
}

func markInProgress(idS string) (string, error) {
	tasks, err := readJSON("tasks.json")
	if err != nil {
		return "", err
	}

	id, err := strconv.Atoi(idS)
	if err != nil {
		return "", err
	}

	found := false

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = StatusInProgress
			newUpdateTask(tasks, i, StatusChange, time.Now())
			found = true
			break
		}
	}

	if !found {
		return "", errors.New("task doesn't exist")
	}

	err = updateJSON("tasks.json", tasks)
	if err != nil {
		return "", err
	}

	return "successfully marked as in progress", nil
}

func markDone(idS string) (string, error) {
	tasks, err := readJSON("tasks.json")
	if err != nil {
		return "", err
	}

	id, err := strconv.Atoi(idS)
	if err != nil {
		return "", err
	}

	found := false

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = StatusDone
			newUpdateTask(tasks, i, StatusChange, time.Now())
			found = true
			break
		}
	}

	if !found {
		return "", errors.New("task doesn't exist")
	}

	err = updateJSON("tasks.json", tasks)
	if err != nil {
		return "", err
	}

	return "successfully marked as done", nil
}

func viewTask(idS string) (string, error) {
	tasks, err := readJSON("tasks.json")
	if err != nil {
		return "", err
	}

	id, err := strconv.Atoi(idS)
	if err != nil {
		return "", err
	}

	found := false

	var result string
	for i := range tasks {
		if tasks[i].ID == id {
			task := tasks[i]

			result = fmt.Sprintln("ID:", task.ID)
			result += fmt.Sprintln("Description:", task.Description)
			result += fmt.Sprintln("Status:", task.Status)
			result += fmt.Sprintln("Created at:", task.CreatedAt)
			result += fmt.Sprintln("Updated at:", task.UpdatedAt)
			result += fmt.Sprintln("History:")
			for j, t := range task.History {
				result += fmt.Sprintln("    ", j, "- Update description:", t.UpdateDesc, "- Update time:", t.UpdateTime)
			}
			found = true
			break
		}
	}

	if !found {
		return "", errors.New("task doesn't exist")
	}

	return result, nil
}

func newUpdateTask(tasks []Task, i int, upDesc UpdateDesc, time time.Time) {
	tasks[i].UpdatedAt = time
	tasks[i].History = append(tasks[i].History, UpdateRecord{UpdateTime: time, UpdateDesc: upDesc})
}

func main() {
	createFile()

	args := os.Args
	if len(args) < 2 {
		printUsage()
		return
	}

	cmd := strings.ToLower(args[1])
	switch cmd {
	case "add":
		if len(args) < 3 {
			fmt.Println("Error: add requires a description")
			return
		}
		desc := strings.Join(args[2:], " ")
		msg, err := addTask(desc)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(msg)

	case "update":
		if len(args) < 4 {
			fmt.Println("Error: update requires <id> <description>")
			return
		}
		id := args[2]
		desc := strings.Join(args[3:], " ")
		msg, err := updateTask(id, desc)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(msg)

	case "delete":
		if len(args) < 3 {
			fmt.Println("Error: delete requires <id>")
			return
		}
		msg, err := deleteTask(args[2])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(msg)

	case "list":
		filter := "all"
		if len(args) >= 3 {
			filter = args[2]
		}
		msg, err := listTasks(filter)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if msg != "" {
			fmt.Println(msg)
		}

	case "mark-in-progress":
		if len(args) < 3 {
			fmt.Println("Error: mark-in-progress requires <id>")
			return
		}
		msg, err := markInProgress(args[2])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(msg)

	case "mark-done":
		if len(args) < 3 {
			fmt.Println("Error: mark-done requires <id>")
			return
		}
		msg, err := markDone(args[2])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(msg)

	case "view-task":
		if len(args) < 3 {
			fmt.Println("Error: view-task requires <id>")
			return
		}
		msg, err := viewTask(args[2])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Print(msg)

	case "repl":
		runRepl()

	case "help", "-h", "--help":
		printUsage()

	default:
		fmt.Println("Error: unknown command:", cmd)
		printUsage()
	}
}

func runRepl() {
	fmt.Println("Task tracker (REPL mode)")
	fmt.Println("Available commands:\n  add <task>\n  update <id> <task>\n  delete <id>\n  list <all|todo|in-progress|done>\n  mark-in-progress <id>\n  mark-done <id>\n  view-task <id>\n  help")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("cli> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		res, err := process(input)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		if strings.TrimSpace(res) != "" {
			fmt.Println(res)
		}
	}
}

func printUsage() {
	fmt.Println("Task Tracker CLI")
	fmt.Println("Usage: task-cli <command> [args]")
	fmt.Println("Commands:")
	fmt.Println("  add <description>")
	fmt.Println("  update <id> <new description>")
	fmt.Println("  delete <id>")
	fmt.Println("  list [all|todo|in-progress|done]")
	fmt.Println("  mark-in-progress <id>")
	fmt.Println("  mark-done <id>")
	fmt.Println("  view-task <id>")
	fmt.Println("  repl    (start interactive mode)")
	fmt.Println("  help    (show this message)")
}

func createFile() {
	_, err := os.Stat("tasks.json")
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	if os.IsNotExist(err) {
		if _, err := os.Create("tasks.json"); err != nil {
			log.Fatal(err)
		}
	}
}

func process(input string) (string, error) {
	tokens := strings.Fields(strings.TrimSpace(input))

	// strings.Fields never returns nil; checking length is sufficient
	if len(tokens) == 0 {
		return "", nil
	}

	action := strings.ToLower(tokens[0])

	switch action {
	case "add":
		if len(tokens) < 2 {
			return "", errors.New("missing task description")
		}
		description := strings.Join(tokens[1:], " ")
		return handleTaskOp(addTask(description))

	case "update":
		if len(tokens) < 3 {
			return "", errors.New("missing task ID or updated description")
		}
		description := strings.Join(tokens[2:], " ")
		return handleTaskOp(updateTask(tokens[1], description))

	case "delete":
		if len(tokens) < 2 {
			return "", errors.New("missing task ID")
		}
		return handleTaskOp(deleteTask(tokens[1]))

	case "list":
		// Allows optional filter (e.g., "list done" or just "list")
		var filter string
		if len(tokens) > 1 {
			filter = tokens[1]
		}
		return handleTaskOp(listTasks(filter))

	case "mark-in-progress":
		if len(tokens) < 2 {
			return "", errors.New("missing task ID")
		}
		return handleTaskOp(markInProgress(tokens[1]))

	case "mark-done":
		if len(tokens) < 2 {
			return "", errors.New("missing task ID")
		}
		return handleTaskOp(markDone(tokens[1]))

	case "view-task":
		if len(tokens) < 2 {
			return "", errors.New("missing task ID")
		}
		return handleTaskOp(viewTask(tokens[1]))

	default:
		return "", errors.New("invalid command")
	}
}

// Helper function to remove boilerplate error checking and logging
func handleTaskOp(res string, err error) (string, error) {
	if err != nil {
		return "", err
	}
	return res, nil
}

func readJSON(filename string) ([]Task, error) {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, errors.New("error reading file")
	}

	var tasks []Task

	if len(fileData) != 0 {
		err = json.Unmarshal(fileData, &tasks)
		if err != nil {
			return nil, errors.New("error parsing json")
		}
	}

	return tasks, nil
}

func updateJSON(filename string, tasks []Task) error {
	updatedJSON, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return errors.New("error marshalling the tasks")
	}

	err = os.WriteFile(filename, updatedJSON, 0644)
	if err != nil {
		return errors.New("error saving the tasks")
	}

	return nil
}

func ParseTaskStatus(s string) (Status, error) {
	// Clean the input to make parsing case-insensitive
	cleaned := strings.ToLower(strings.TrimSpace(s))

	switch cleaned {
	case string(StatusTodo):
		return StatusTodo, nil
	case string(StatusInProgress):
		return StatusInProgress, nil
	case string(StatusDone):
		return StatusDone, nil
	default:
		// Return an error if the string is unrecognized
		return "", fmt.Errorf("invalid status %q: must be 'todo', 'in-progress', or 'done'", s)
	}
}

func nextID(tasks []Task) int {
	mx := 0
	for _, t := range tasks {
		mx = max(t.ID, mx)
	}

	return mx
}
