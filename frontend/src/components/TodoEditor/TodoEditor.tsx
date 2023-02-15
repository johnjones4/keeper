import React, { useEffect, useState } from 'react';
import { Todos } from '../../models/Todos';
import './TodoEditor.css'

interface TodoEditorProps {
  value: string
  onChange(value: string): void
}

const TodoEditor = (props: TodoEditorProps) => {
  const [todos, setTodos] = useState(new Todos())
  const [newTodoText, setNewTodoText] = useState('')

  useEffect(() => {
    setTodos(Todos.fromSource(props.value))
  }, [props.value])

  const todoChanged = (index: number, checked: boolean) => {
    const newTodos = todos.toggleItem(index, checked)
    setTodos(newTodos)
    props.onChange(newTodos.text())
  }

  const commitNewTodo = () => {
    if (newTodoText.trim() === '') {
      return
    }
    const newTodos = todos.addItem(newTodoText)
    setTodos(newTodos)
    setNewTodoText('')
    props.onChange(newTodos.text())
  }

  const removeTodo = (index: number) => {
    const newTodos = todos.removeItem(index)
    setTodos(newTodos)
    props.onChange(newTodos.text())
  }

  const indent = (index: number, direction: number) => {
    const newTodos = todos.indentItem(index, direction)
    setTodos(newTodos)
    props.onChange(newTodos.text())
  }

  return (
    <div className='TodoEditor'>
      { todos.todos.map((todo, i) => {
        return (
          <div key={todo.text} className={`Todo ${todo.done ? 'Todo-Done' : 'Todo-Undone'}`} style={{paddingLeft: `${todo.indent}em`}}>
            <label>
              <input type="checkbox" checked={todo.done} onChange={(e) => todoChanged(i, e.target.checked)} />
              <span className='Todo-Text'>{ todo.text }</span>
            </label>
            <button onClick={() => indent(i, -1)} disabled={todo.indent === 0}>&larr;</button>
            <button onClick={() => indent(i, 1)}>&rarr;</button>
            <button onClick={() => removeTodo(i)}>x</button>
          </div>
        )
      }) }
      <form onSubmit={(e) => {
        e.preventDefault()
        commitNewTodo()
      }}>
        <div className='TodoNew InputGroup'>
          <input type='text' className='Input' value={newTodoText} onChange={(e) => setNewTodoText(e.target.value)} />
          <button className='Btn' onClick={() => commitNewTodo()}>+</button>
        </div>
      </form>
    </div>
  )
}

export default TodoEditor
