export default class NoteNode {
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
