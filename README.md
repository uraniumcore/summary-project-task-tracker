# Task Tracker CLI

https://roadmap.sh/projects/task-tracker

A simple command‑line app to track your tasks with JSON persistence. Supports adding, updating, deleting, viewing, listing (with filters), and changing task status. Data is stored in a local tasks.json file in the current directory.

This version supports positional CLI arguments (single‑shot commands) and also includes an optional interactive REPL mode.

## Build

Prerequisites: Go 1.21+ (any recent Go should work)

```bash
# from the project root
go build -o task-cli
```

## Usage (positional CLI)

```bash
# Adding a new task
./task-cli add "Buy groceries"

# Updating a task
./task-cli update 1 "Buy groceries and cook dinner"

# Deleting a task
./task-cli delete 1

# Marking a task as in progress or done
./task-cli mark-in-progress 2
./task-cli mark-done 2

# Listing all tasks or by status
./task-cli list            # same as "list all"
./task-cli list all
./task-cli list todo
./task-cli list in-progress
./task-cli list done

# Viewing a single task with full details and history
./task-cli view-task 3

# Help
./task-cli help
```

## Optional REPL mode
If you prefer an interactive shell, start the REPL:

```bash
./task-cli repl
```
Then type commands like:

```
add Read a book
update 1 Read two books
delete 1
list todo
mark-in-progress 2
mark-done 2
view-task 2
```

## Data model
Each task is stored with the following properties:
- id (int)
- description (string)
- status (string): one of "todo", "in-progress", "done"
- createdAt (RFC3339 timestamp)
- updatedAt (RFC3339 timestamp)
- history (array of update records), where each record contains:
  - update_desc (enum-like int): StatusChange | TaskChange
  - update_time (timestamp)

## Storage
- File: tasks.json in the current working directory
- The file is created automatically if it does not exist

## Notes
- IDs are allocated monotonically based on the maximum existing ID in the file (next = maxID + 1).
- The CLI returns user‑friendly errors (e.g., invalid IDs, unknown commands, or invalid statuses).
