import React, { useEffect, useState } from 'react';
import { Note } from './Note';

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
      <textarea value={noteBody} onChange={(event) => bodyChange(event.target.value)} />
    </div>
  )
}

export default NoteDetail
