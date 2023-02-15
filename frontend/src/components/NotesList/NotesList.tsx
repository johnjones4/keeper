import React, { useEffect, useState } from 'react';
import AddNote from '../AddNote/AddNote';
import { Note } from '../../models/Note';
import './NotesList.css'
import NoteNode from '../../models/NoteNode';

interface NotesListProps {
  onNoteSelected(key: string): void
  onNewNote(prefix?: string): void
  onSearch(query: string): void
  onWantsSubdirectory(dir: string): void
  onDiscardSubdirectory(dir: string): void
  notes: string[]
}

const NotesList = (props: NotesListProps) => {
  const [tree, setTree] = useState(new NoteNode("", undefined))
  const [query, setQuery] = useState('')

  useEffect(() => {
    const localTree = new NoteNode("", undefined)
    props.notes.forEach(note => localTree.addNote(note))
    setTree(localTree)
  }, [props.notes])

  const renderTree = (node: NoteNode): any => {
    return (
      <li key={node.name}>
        <div className='NoteTitle'>
          <button onClick={() => node.children.length > 0 || node.notes.length > 0 ? props.onDiscardSubdirectory(node.getBasePath()) : props.onWantsSubdirectory(node.getBasePath())} className='NoteTitleText'>{node.name === '' ? 'Root' : node.name}</button>
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
    <div className='NotesList'>
      <div className='NotesListSearch InputGroup'>
        { query !== '' && (
          <button className='Btn' onClick={() => {
            props.onSearch('')
            setQuery('')
          }}>X</button>
        ) }
        <input className='Input' value={query} onChange={(e) => setQuery(e.target.value)} />
        <button className='Btn' onClick={() => props.onSearch(query)}>Search</button>
      </div>
      <ul>
        { renderTree(tree) }
      </ul>
    </div>
  )
}

export default NotesList
