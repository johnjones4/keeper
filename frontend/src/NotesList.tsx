import React, { useEffect, useState } from 'react';
import AddNote from './AddNote';
import { Note } from './Note';

interface NotesListProps {
  onNoteSelected(key: string): void
  onNewNote(prefix: string): void
  notes: string[]
}

class NoteNode {
  parent?: NoteNode
  name: string
  children: NoteNode[]
  notes: string[]
  
  constructor(name: string, parent?: NoteNode) {
    this.parent = parent
    this.name = name
    this.children = []
    this.notes = []
  }

  addNote(key: string) {
    const allParts = key.split('/')
    if (allParts.length <= 1 || allParts[0] !== "") {
      throw Error('Unexpected key: ' + key)
    }
    this._addNote(allParts.slice(1))
  }

  _addNote(parts: string[]) {
    if (parts.length === 1) {
      this.notes.push(parts[0])
    } else {
      const root = parts.shift()
      let child = this.children.find(c => c.name === root)
      if (!child) {
        child = new NoteNode(root as string, this)
        this.children.push(child)
      }
      child._addNote(parts)
    }
  }

  getPath(noteName: string): string {
    return this.getBasePath() + '/' + noteName
  }

  getBasePath(): string {
    if (this.parent) {
      return this.parent.getBasePath() + '/' + this.name
    } else {
      return this.name
    }
  }
}

const NotesList = (props: NotesListProps) => {
  const [tree, setTree] = useState(new NoteNode("", undefined))

  useEffect(() => {
    const localTree = new NoteNode("", undefined)
    props.notes.forEach(note => localTree.addNote(note))
    setTree(localTree)
  }, [props.notes])

  const renderTree = (node: NoteNode): any => {
    return (
      <li key={node.name}>
        <div className='NoteTitle'>
          <span className='NoteTitleText'>{node.name === '' ? 'Root' : node.name}</span>
          <button onClick={() => props.onNewNote(node.getBasePath())} className='ButtonAddNote'>+</button>
        </div>
        <ul>
          { node.children.map(c => renderTree(c)) }
        </ul>
        <ul>
          { node.notes.map(n => (
            <li key={n}>
              <button onClick={() => props.onNoteSelected(node.getPath(n))} className='ButtonOpenNote'>
                {n}
              </button>
            </li>
          )) }
        </ul>
      </li>
    )
  }

  return (
    <div className='NoteList'>
      <ul>
        { renderTree(tree) }
      </ul>
    </div>
  )
}

export default NotesList
