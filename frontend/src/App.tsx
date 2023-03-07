import React, { useEffect, useState } from 'react';
import AddNote from './components/AddNote/AddNote';
import './App.css';
import Login from './components/Login';
import { Message, MessageType } from './models/Message';
import { Note } from './models/Note';
import NoteDetail from './components/NoteDetail/NoteDetail';
import NotesList from './components/NotesList/NotesList';
import tokenManager from './models/TokenManager'
import NoteSet from './models/NoteSet';


function App() {
  const [loggedIn, setLoggedIn] = useState(tokenManager.isReady())
  const [notes, setNotes] = useState(new NoteSet([]))
  const [note, setNote] = useState(null as string|null)
  const [message, setMessage] = useState(null as null|Message)
  const [addNotePrefix, setAddNotePrefix] = useState(null as string|null)

  const loadNotesInitial = async () => {
    try {
      const resp = await Note.getNotesDir('/')
      setNotes(new NoteSet(resp.notes))
    } catch (e) {
      showMessage({type: MessageType.error, message: `${e}`})
    }
  }

  const loadNotesDir = async (dir: string) => {
    try {
      const resp = await Note.getNotesDir(dir)
      setNotes(notes.update(resp.notes))
    } catch (e) {
      showMessage({type: MessageType.error, message: `${e}`})
    }
  }

  const removeNotesDir = (dir: string) => {
    setNotes(new NoteSet(notes.notes.filter(note => {
      console.log(note, dir)
      return note === dir || !note.startsWith(dir)
    })))
  }

  const doSearch = async (query: string) => {
    if (query === '') {
      loadNotesInitial()
      return
    }
    const localNotes = await Note.search(query)
    setNotes(new NoteSet(localNotes.notes))
  }

  useEffect(() => {
    loadNotesInitial()
  }, [loggedIn])

  const showMessage = (m: Message) => {
    setMessage(m)
    setTimeout(() => setMessage(null), 5000)
  }

  const handleError = (e: any) => {
    if (`${e}` === 'Error: needs login') {
      setLoggedIn(false)
      return
    }
    showMessage({type: MessageType.error, message: `${e}`})
  }

  if (!loggedIn) {
    return (<Login 
      onError={e => showMessage({type: MessageType.error, message: `${e}`})}
      onLogin={() => setLoggedIn(true)}
    />)
  }

  return (
    <div className='NoteApp'>
      <NotesList 
        notes={notes.notes}
        onNoteSelected={k => setNote(k)} 
        onNewNote={prefix => {
          if (prefix) {
            loadNotesDir(prefix)
          }
          setAddNotePrefix(prefix !== undefined ? prefix : null)
        }}
        onSearch={q => doSearch(q)}
        onWantsSubdirectory={d => loadNotesDir(d)}
        onDiscardSubdirectory={d => removeNotesDir(d)}
        />
      { note && (<NoteDetail 
        noteKey={note} 
        onError={e => handleError(e)}
        onMessage={m => showMessage({type: MessageType.alert, message: m})}
        />
      ) }
      { addNotePrefix !== null && (
        <AddNote 
          initialPrefix={addNotePrefix} 
          onError={e =>  handleError(e)}
          onNewNote={n => {
            setNote(n.key)
            setAddNotePrefix(null)
            setNotes(notes.update([n.key]))
          }} 
          onCancel={() => setAddNotePrefix(null)} 
        />
      )}
      { message && <div className={`Message Message-${message.type}`}>{message.message}</div> }
    </div>
  );
}

export default App;
