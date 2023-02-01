import React, { useEffect, useState } from 'react';
import AddNote from './AddNote';
import './App.css';
import Login from './Login';
import { Message, MessageType } from './Message';
import { Note } from './Note';
import NoteDetail from './NoteDetail';
import NotesList from './NotesList';
import tokenManager from './TokenManager'


function App() {
  const [loggedIn, setLoggedIn] = useState(tokenManager.isReady())
  const [notes, setNotes] = useState([] as string[])
  const [note, setNote] = useState(null as string|null)
  const [message, setMessage] = useState(null as null|Message)
  const [addNotePrefix, setAddNotePrefix] = useState(null as string|null)

  const loadNotes = async () => {
    try {
      let endOfList = false
      let localNotes = [] as string[]
      let page = ''
      while (!endOfList) {
        const resp = await Note.getNotes(page)
        page = resp.nextPage
        endOfList = page === ''
        if (resp.notes.length > 0) {
          localNotes = localNotes.concat(resp.notes)
          setNotes(localNotes)
        }
      }
    } catch (e) {
      showMessage({type: MessageType.error, message: `${e}`})
    }
  }

  const doSearch = async (query: string) => {
    if (query === '') {
      loadNotes()
      return
    }
    const localNotes = await Note.search(query)
    setNotes(localNotes.notes)
  }

  useEffect(() => {
    loadNotes()
  }, [loggedIn])

  const showMessage = (m: Message) => {
    setMessage(m)
    setTimeout(() => setMessage(null), 5000)
  }

  const handleError = (e: any) => {
    if (`${e}` === 'needs login') {
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
        notes={notes}
        onNoteSelected={k => setNote(k)} 
        onNewNote={prefix => setAddNotePrefix(prefix !== undefined ? prefix : null)}
        onSearch={q => doSearch(q)}
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
            setNotes(notes.concat(n.key))
          }} 
          onCancel={() => setAddNotePrefix(null)} 
        />
      )}
      { message && <div className={`Message Message-${message.type}`}>{message.message}</div> }
    </div>
  );
}

export default App;
