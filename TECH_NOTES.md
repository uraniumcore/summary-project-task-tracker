# Technical Notes: What this project uses

This document highlights the Go language features and standard library packages used in this project, from basics to the more advanced pieces that were the trickiest.

## Basics
- for loops: used in many places to iterate over tasks (e.g., scanning to find a task by ID, computing nextID)
- fmt: for printing user messages and formatted output in the CLI
- strings: splitting user input (Fields/Join), trimming whitespace, lowercasing commands
- errors: creating and returning meaningful errors from operations

## Files and JSON
- os: reading/writing files, checking/creating tasks.json on startup
- encoding/json: marshalling/unmarshalling the tasks to and from JSON with struct tags and pretty printing (MarshalIndent)

## Time handling
- time: timestamps for createdAt and updatedAt, plus recording update history moments

## Data modeling with structs and embedded/nested data
- Structs: Task and UpdateRecord model your domain data cleanly
- Nested/embedded-like design: Task includes a History slice of UpdateRecord values, which acts like an embedded collection of changes per task

## Enum-like constants
- Status is a custom string type with constants (todo, in-progress, done). This avoids hard-coded strings throughout the code.
- UpdateDesc is an int-based enum with iota to differentiate change kinds (StatusChange vs TaskChange), with a String method for readable output.

## CLI parsing approaches
- Positional CLI mode: Uses os.Args to accept commands of the form `task-cli add "text"`, etc.
- REPL mode: A simple interactive loop using bufio.Reader + process() that parses commands line-by-line.

## Error handling improvements
- Validate and parse IDs (strconv.Atoi) with proper error checks
- Detect missing/unknown task IDs and return clear errors
- Centralized status parsing with ParseTaskStatus to ensure only valid values are accepted

## ID allocation strategy
- nextID computes the maximum existing ID and new tasks use maxID+1. This keeps IDs unique and monotonically increasing, even if tasks are deleted or re-ordered.

## Output consistency
- viewTask returns a single formatted string for the caller to print (consistent with other command functions)
- listTasks prints details to stdout and returns a short summary for the caller

## Possible future enhancements
- Persist a nextId counter alongside tasks for O(1) ID allocation
- Add `mark-todo <id>` and a `help` subcommand with richer docs
- Sort list output (e.g., by id ascending) and format timestamps with RFC3339 consistently in all output
- Add tests for core behaviors (add/update/delete/list/status changes)
