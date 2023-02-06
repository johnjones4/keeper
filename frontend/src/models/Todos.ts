export interface Todo {
  done: boolean
  text: string
}

enum TodoState {
  undone = '[ ] ',
  done = '[x] '
}

export class Todos {
  todos: Todo[]

  constructor() {
    this.todos = []
  }

  static fromSource(source: string) {
    const todos = new Todos()
    todos.todos = source.split('\n').map(line => {
      if (line.trim() === '') {
        return null
      }
      return {
        done: line.substring(0, 4) !== TodoState.undone,
        text: line.substring(4)
      }
    }).filter(i => !!i) as Todo[]
    return todos
  }

  static fromTodos(todos: Todo[]) {
    const t = new Todos()
    t.todos = todos
    return t
  }

  text(): string {
    return this.todos.map(todo => {
      return (todo.done ? TodoState.done : TodoState.undone) + todo.text
    }).join('\n')
  }

  toggleItem(index: number, value: boolean): Todos {
    return Todos.fromTodos(this.todos.map((todo, i) => {
      if (i === index) {
        todo.done = value
      }
      return todo
    }))
  }

  removeItem(index: number): Todos {
    return Todos.fromTodos(this.todos.filter((_, i) => i !== index))
  }

  addItem(text: string): Todos {
    return Todos.fromTodos(this.todos.concat([{
      done: false,
      text
    }]))
  }
}
