package main

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type runtime struct {
	list        *todoList
	screen      tcell.Screen
	selected    int
	lastDelete  time.Time
	inputText   string
	inputCursor int
	moving      bool
}

const (
	doubleDelete = time.Second / 2
)

func newRuntime(list *todoList) *runtime {
	return &runtime{
		list:     list,
		selected: -1,
	}
}

func (rt *runtime) render() {
	w, h := rt.screen.Size()
	rt.screen.Clear()

	listHeight := h - 2

	listIndex := 0
	if rt.selected >= listHeight {
		listIndex = rt.selected - listHeight + 1
	}
	y := 0
	for y < listHeight && listIndex < len(rt.list.todos) {
		todo := rt.list.todos[listIndex]
		lineStr := todo.string()
		prefixIndent := 0
		if rt.moving && listIndex == rt.selected {
			lineStr = ">" + lineStr
			prefixIndent = 1
		}
		x := 0
		for i, c := range lineStr {
			var comb []rune
			w1 := runewidth.RuneWidth(c)
			if w1 == 0 {
				comb = []rune{c}
				c = ' '
				w = 1
			}
			style := tcell.StyleDefault.
				Underline(!rt.moving && i == todo.indent+1+prefixIndent && listIndex == rt.selected).
				StrikeThrough(todo.done && i > todo.indent+3+prefixIndent)
			rt.screen.SetContent(x, y, c, comb, style)
			x += w1
			if x >= w {
				break
			}
		}
		y++
		listIndex++
	}

	x := 0
	start := 0
	text := rt.inputText + " "
	if rt.inputCursor >= w && rt.selected < 0 {
		start = rt.inputCursor - w + 1
	}
	for i, c := range text[start:] {
		var comb []rune
		w1 := runewidth.RuneWidth(c)
		if w1 == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		style := tcell.StyleDefault.
			Reverse(start+i == rt.inputCursor && rt.selected < 0)
		rt.screen.SetContent(x, h-1, c, comb, style)
		x += w1
		if x >= w {
			break
		}
	}
	rt.screen.Sync()
}

func (rt *runtime) handleKeyForTodos(event *tcell.EventKey) error {
	switch event.Key() {
	case tcell.KeyTAB:
		rt.selected = -1
		rt.moving = false
		rt.render()
	case tcell.KeyLeft:
		if rt.list.todos[rt.selected].indent == 0 {
			return nil
		}
		rt.list.todos[rt.selected].indent--
		err := rt.list.save()
		if err != nil {
			return err
		}
		rt.render()
	case tcell.KeyRight:
		rt.list.todos[rt.selected].indent++
		err := rt.list.save()
		if err != nil {
			return err
		}
		rt.render()
	case tcell.KeyDown:
		rt.lastDelete = time.Time{}
		if rt.selected < len(rt.list.todos)-1 {
			if rt.moving {
				temp := rt.list.todos[rt.selected+1]
				rt.list.todos[rt.selected+1] = rt.list.todos[rt.selected]
				rt.list.todos[rt.selected] = temp
			}
			rt.selected++
		} else {
			rt.selected = -1
			rt.moving = false
		}
		rt.render()
	case tcell.KeyUp:
		rt.lastDelete = time.Time{}
		if rt.selected > 0 {
			if rt.moving {
				temp := rt.list.todos[rt.selected-1]
				rt.list.todos[rt.selected-1] = rt.list.todos[rt.selected]
				rt.list.todos[rt.selected] = temp
			}
			rt.selected--
		} else {
			rt.selected = -1
			rt.moving = false
		}
		rt.render()
	case tcell.KeyEnter:
		if rt.moving {
			rt.moving = false
		} else {
			rt.lastDelete = time.Time{}
			err := rt.list.toggle(rt.selected)
			if err != nil {
				return err
			}
		}
		rt.render()
	case tcell.KeyDEL, tcell.KeyDelete:
		if rt.lastDelete.IsZero() {
			rt.lastDelete = time.Now()
		} else if time.Now().Before(rt.lastDelete.Add(doubleDelete)) {
			rt.lastDelete = time.Time{}
			newSelected := rt.selected
			if rt.selected == len(rt.list.todos)-1 {
				newSelected--
			}
			err := rt.list.remove(rt.selected)
			if err != nil {
				return err
			}
			rt.selected = newSelected
			rt.render()
		}
	case tcell.KeyRune:
		r := event.Rune()
		if string(r) == " " {
			rt.lastDelete = time.Time{}
			err := rt.list.toggle(rt.selected)
			if err != nil {
				return err
			}
			rt.render()
		} else if string(r) == "m" {
			rt.moving = true
			rt.render()
		}
	}
	return nil
}

func (rt *runtime) handleKeyForInput(event *tcell.EventKey) error {
	switch event.Key() {
	case tcell.KeyUp:
		rt.selected = len(rt.list.todos) - 1
		rt.render()
	case tcell.KeyTAB, tcell.KeyDown:
		rt.selected = 0
		rt.render()
	case tcell.KeyLeft:
		if rt.inputCursor > 0 {
			rt.inputCursor--
			rt.render()
		}
	case tcell.KeyRight:
		if rt.inputCursor < len(rt.inputText) {
			rt.inputCursor++
			rt.render()
		}
	case tcell.KeyDEL:
		if len(rt.inputText) <= 1 {
			rt.inputText = ""
			rt.inputCursor = 0
		} else {
			newText := make([]byte, len(rt.inputText)-1)
			for i := 0; i < len(newText); i++ {
				if i < rt.inputCursor {
					newText[i] = rt.inputText[i]
				} else {
					newText[i] = rt.inputText[i+1]
				}
			}
			rt.inputText = string(newText)
			if rt.inputCursor > 0 {
				rt.inputCursor--
			}
		}
		rt.render()
	case tcell.KeyEnter:
		err := rt.list.add(rt.inputText)
		if err != nil {
			return err
		}
		rt.inputText = ""
		rt.inputCursor = 0
		rt.render()
	case tcell.KeyRune:
		r := event.Rune()
		rt.render()
		if rt.inputCursor >= len(rt.inputText) {
			rt.inputText += string(r)
		} else {
			newText := make([]byte, len(rt.inputText)+1)
			for i := 0; i < len(newText); i++ {
				if i <= rt.inputCursor {
					newText[i] = rt.inputText[i]
				} else if i == rt.inputCursor+1 {
					newText[i] = byte(r)
				} else {
					newText[i] = rt.inputText[i-1]
				}
			}
			rt.inputText = string(newText)
		}
		rt.inputCursor++
		rt.render()
	}
	return nil
}

func (rt *runtime) run() error {
	var err error
	rt.screen, err = tcell.NewScreen()
	if err != nil {
		return err
	}

	err = rt.screen.Init()
	if err != nil {
		return err
	}

	rt.screen.SetStyle(tcell.StyleDefault)
	rt.screen.Clear()
	rt.render()

	for {
		ev := rt.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				rt.screen.Fini()
				return nil
			} else if rt.selected >= 0 {
				err := rt.handleKeyForTodos(ev)
				if err != nil {
					rt.screen.Fini()
					return err
				}
			} else {
				err := rt.handleKeyForInput(ev)
				if err != nil {
					rt.screen.Fini()
					return err
				}
			}

		case *tcell.EventResize:
			rt.render()
		}
	}
}
