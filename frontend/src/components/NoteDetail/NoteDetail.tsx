import React, { useEffect, useState } from 'react';
import { Note } from '../../models/Note';
import MDEditor from '@uiw/react-md-editor';
import './NoteDetail.css'
import TodoEditor from '../TodoEditor/TodoEditor';
const path = require('path-browserify');

interface NotesListProps {
  noteKey: string
  onError(err: any): void
  onMessage(msg: string): void
}

let timeout: NodeJS.Timeout | null

const NoteDetail = (props: NotesListProps) => {
  const [noteBody, setNoteBody] = useState("")
  const [note, setNote] = useState(null as null|Note)

  const loadNote = async (key: string) => {
    try {
      const note = await Note.getNote(key)
      setNote(note)
      setNoteBody(note.body)
    } catch (e) {
      props.onError(e)
    }
  }

  useEffect(() => {
    loadNote(props.noteKey)
  }, [props.noteKey])

  const bodyChange = (body: string) => {
    setNoteBody(body)
    if (note) {
      note.body = body
    }
    if (timeout) {
      clearTimeout(timeout)
    }
    timeout = setTimeout(async () => {
      if (note) {
        try {
          await note.save()
          props.onMessage('Note saved')
        } catch (e) {
          props.onError(e)
        }
      }
    }, 1000)
  }

  const closeNote = () => {
    setNote(null)
    setNoteBody('')
  }

  const renderEditor = (): any => {
    if (note) {
      const ext = path.extname(note.key)
      switch (ext) {
        case '.md':
          return (
            <MDEditor
              value={noteBody}
              onChange={v => bodyChange(v ? v : '')}
              preview='edit'
            />
          )
        case '.todo':
          return (
            <TodoEditor
              value={noteBody}
              onChange={v => bodyChange(v ? v : '')}
            />
          )
        default:
          return (<textarea className='TexteditorPlain' value={noteBody} onChange={(event) => bodyChange(event.target.value)} />)
      }
    }
  }

  if (!note) {
    return null
  }

  return (
    <div className='NoteDetail'>
      <div className='NoteDetailHead'>
        <div className='NoteDetailHeadTitle'>
          {note.key}
        </div>
        <button className='Btn NoteDetailCloseButton' onClick={() => closeNote()}>Close</button>
      </div>
      { renderEditor() }
    </div>
  )
}

export default NoteDetail
