package main

import (
	"errors"
	"os"
	"strings"
)

type todo struct {
	text   string
	done   bool
	indent int
}

func (t todo) string() string {
	prefix := undoneStr
	if t.done {
		prefix = doneStr
	}
	return makeIndent(t.indent) + prefix + t.text
}

type todoList struct {
	todos []todo
	path  string
}

const (
	doneStr   = "[x] "
	undoneStr = "[ ] "
)

var (
	errorOutOfBounds = errors.New("index out of bounds")
	errorInvalidLine = errors.New("invalid line")
)

func (tl *todoList) open() error {
	contents, err := os.ReadFile(tl.path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if errors.Is(err, os.ErrNotExist) {
		tl.todos = make([]todo, 0)
		return nil
	}

	lines := strings.Split(string(contents), "\n")
	tl.todos = make([]todo, 0)
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := 0
		if line[0] == ' ' {
			for indent < len(line) && line[indent] == ' ' {
				indent++
			}
		} else {
			indent = 0
		}
		if line[indent] == '[' {
			tl.todos = append(tl.todos, todo{
				done:   strings.Index(line[indent:], doneStr) == 0,
				text:   line[indent+4:],
				indent: indent,
			})
		} else {
			return errorInvalidLine
		}
	}

	return nil
}

func makeIndent(i int) string {
	strs := make([]byte, i)
	for j := 0; j < i; j++ {
		strs[j] = ' '
	}
	return string(strs)
}

func (tl *todoList) save() error {
	lines := make([]string, len(tl.todos))
	for i, todo := range tl.todos {
		lines[i] = todo.string()
	}

	contents := strings.Join(lines, "\n")

	err := os.WriteFile(tl.path, []byte(contents), 0777)
	if err != nil {
		return err
	}

	return nil
}

func (tl *todoList) toggle(i int) error {
	if i < 0 || i >= len(tl.todos) {
		return errorOutOfBounds
	}
	tl.todos[i].done = !tl.todos[i].done
	return tl.save()
}

func (tl *todoList) remove(index int) error {
	if index < 0 || index >= len(tl.todos) {
		return errorOutOfBounds
	}
	if len(tl.todos) == 1 {
		tl.todos = make([]todo, 0)
		return nil
	}

	newList := make([]todo, len(tl.todos)-1)
	for i := 0; i < len(newList); i++ {
		if i < index {
			newList[i] = tl.todos[i]
		} else {
			newList[i] = tl.todos[i+1]
		}
	}
	tl.todos = newList

	return tl.save()
}

func (tl *todoList) add(text string) error {
	tl.todos = append(tl.todos, todo{text: text})
	return tl.save()
}
