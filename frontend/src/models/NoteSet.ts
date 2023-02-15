import { Note } from "./Note"

export default class NoteSet {
  notes: string[]
  
  constructor(items: string[]) {
    this.notes = items
  }

  update(items: string[]): NoteSet {
    return new NoteSet(this.notes.concat(items.filter(note => {
      return this.notes.findIndex(note1 => note === note1) < 0
    })))
  }
}
