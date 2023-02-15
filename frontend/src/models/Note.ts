import { Base64 } from 'js-base64'
import tokenManager from './TokenManager'

interface NoteResponse {
  key: string
  body: string
  modified: string
}

export interface NotesResponse {
  notes: string[]
  nextPage: string
}

const errorHandler = (res: Response, json: any) => {
  if (res.status !== 200) {
    if (res.headers.has('X-Show-Login')) {
      throw new Error('needs login')
    }
    throw new Error(json.message)
  }
}

export class Note {
  key: string
  body: string
  modified: Date

  private constructor(response: NoteResponse) {
    this.key = response.key
    this.body = response.body
    this.modified = new Date(Date.parse(response.modified))
  }

  static async getNotes(page: string): Promise<NotesResponse> {
    const res = await fetch(`/api/note?page=${encodeURIComponent(page)}`, {
      headers: tokenManager.getHeaders()
    })
    const json = await res.json()
    errorHandler(res, json)
    return json as NotesResponse
  }

  static async getNotesDir(dir: string): Promise<NotesResponse> {
    const res = await fetch(`/api/note?dir=${encodeURIComponent(dir)}`, {
      headers: tokenManager.getHeaders()
    })
    const json = await res.json()
    errorHandler(res, json)
    return json as NotesResponse
  }

  static async getNote(key: string): Promise<Note> {
    const id = Base64.encode(key)
    const res = await fetch(`/api/note/${encodeURIComponent(id)}`, {
      headers: tokenManager.getHeaders()
    })
    const json = await res.json()
    errorHandler(res, json)
    return new Note(json as NoteResponse)
  }

  static async search(query: string): Promise<NotesResponse> {
    const res = await fetch(`/api/note?q=${encodeURIComponent(query)}`, {
      headers: tokenManager.getHeaders()
    })
    const json = await res.json()
    errorHandler(res, json)
    return json as NotesResponse
  }

  static async newNote(key: string, body: string): Promise<Note> {
    const res = await fetch('/api/note', {
      method: 'POST',
      body: JSON.stringify({key, body}),
      headers: tokenManager.getHeaders()
    })
    const json = await res.json()
    errorHandler(res, json)
    return new Note(json as NoteResponse)
  }

  async save(): Promise<Note> {
    const id = Base64.encode(this.key)
    const res = await fetch(`/api/note/${encodeURIComponent(id)}`, {
      method: 'PUT',
      body: JSON.stringify(this),
      headers: tokenManager.getHeaders()
    })
    const json = await res.json()
    errorHandler(res, json)
    return new Note(json as NoteResponse)
  }
}
