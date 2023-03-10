import React, { useState } from 'react';
import { Note } from '../../models/Note';
import sanitize from 'sanitize-filename';
import './AddNote.css'
const path = require('path-browserify');

interface AddNoteProps {
  initialPrefix: string
  onNewNote(note: Note): void
  onError(err: any): void
  onCancel(): void
}

const AddNote = (props: AddNoteProps) => {
  const [prefix, setPrefix] = useState(props.initialPrefix)
  const [name, setName] = useState('')
  const [format, setFormat] = useState('txt')

  const newNote = async () => {
    if (!isValidName()) {
      return 
    }
    const cleanPrefix = path.normalize(path.isAbsolute(prefix) ? prefix : path.join('/', prefix))
    const key = path.normalize(path.join(cleanPrefix, name + '.' + format));
    console.log(key)
    try {
      const note = await Note.newNote(key, format !== 'todo' ? 'New Note' : '')
      props.onNewNote(note)
    } catch (e) {
      props.onError(e)
    }
  }

  const setPrefixSafe = (text: string) => {
    setPrefix(text.split('/').map(d => sanitize(d.trim())).join('/'))
  }

  const setNameSafe = (text: string) => {
    setName(sanitize(text.trim()))
  }

  const isValidName = (): boolean => {
    return name !== ''
  }

  return (
    <div className='AddNote'>
      <div className='AddNoteHead'>Add Note</div>
      <div className='InputField'>
        <label htmlFor='prefix'>Prefix</label>
        <input className='Input' value={prefix} onChange={e => setPrefixSafe(e.target.value)} name='prefix' />
      </div>
      <div className='InputField'>
        <label htmlFor='name'>Name</label>
        <input className='Input' value={name} onChange={e => setNameSafe(e.target.value)} name='name' />
      </div>
      <div className='InputField'>
        <label htmlFor='format'>Format</label>
        <select className='Input' value={format} onChange={e => setFormat(e.target.value)} name='format'>
          <option value="txt">Text</option>
          <option value="md">Markdown</option>
          <option value="todo">Todo</option>
        </select>
      </div>
      <div className='InputFooter'>
        <button className='Btn AddNoteCancel' onClick={() => props.onCancel()}>Cancel</button>
        <button className='Btn AddNoteOk' disabled={!isValidName()} onClick={() => newNote()}>Save</button>
      </div>
    </div>
  )
}

export default AddNote
