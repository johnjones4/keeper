import { defineStore } from 'pinia';
import { Note, Notes, Message } from 'src/client/types.gen'
import { getNote, getNoteByKey, postNote, putNoteByKey } from 'src/client/services.gen'
import { useErrorStore } from './error-store';
import { useTokenStore } from './token-store';

export interface NoteNode {
  name: string;
  key: string;
}

export interface DirNode {
  name: string;
  fullpath: string;
  notes: NoteNode[];
  dirs: DirNode[];
}

export const useNotesStore = defineStore('notes', {
  state: () => ({
    notes: {
      name: '',
      fullpath: '/',
      notes: [] as NoteNode[],
      dirs: [] as DirNode[],
    } as DirNode,
    note: {
      key: '',
      body: ''
    },
  }),
  actions: {
    clear() {
      this.note = {
        key: '',
        body: ''
      }
      this.notes = {
        name: '',
        fullpath: '/',
        notes: [] as NoteNode[],
        dirs: [] as DirNode[],
      } as DirNode
    },
    async loadRootDirectory() {
      await this.loadDirectory('/')
    },
    async loadDirectory(dir: string) {
      try {
        const response = await getNote({
          dir,
          authorization: useTokenStore().authorization,
        })
        if ((response as Message).ok !== undefined) {
          throw new Error((response as Message).message)
        }
        const notes = (response as Notes).notes;
        const dirPath = dir === '/' ? [] : dir.split('/')
        let currentNode = this.notes
        dirPath.forEach(subDir => {
          if (!subDir) {
            return
          }
          if (!currentNode) {
            return
          }
          const nextDir = currentNode.dirs.find(d => d.name === subDir);
          if (!nextDir) {
            const newDir = {
              name: subDir,
              fullpath: currentNode.fullpath + '/' + subDir,
              notes: [] as NoteNode[],
              dirs: [] as DirNode[],
            }
            currentNode.dirs.push(newDir);
            currentNode = newDir;
          } else {
            currentNode = nextDir;
          }
        });
        // currentNode.notes = [];
        // currentNode.dirs = [];
        notes.forEach(n => {
          const basename = n.split('/').reverse()[0];
          if (n.indexOf('.') >= 0) {
            const key = btoa(n);
            const found = currentNode.notes.find(nn => key === nn.key);
            if (!found) {
              currentNode.notes.push({
                name: basename,
                key,
              })
            }
          } else {
            const found = currentNode.dirs.find(dd => n === dd.fullpath);
            if (!found) {
              currentNode.dirs.push({
                name: basename,
                fullpath: n,
                notes: [] as NoteNode[],
                dirs: [] as DirNode[],
              })
            }
          }
        })
      } catch (err) {
        useErrorStore().set(new Error(`${err}`));
      }
    },
    async selectNote(key: string) {
      try {
        const response = await getNoteByKey({
          key: key,
          authorization: useTokenStore().authorization,
        })
        if ((response as Message).ok !== undefined) {
          throw new Error((response as Message).message)
        }
        this.note = (response as Note)
      } catch (err) {
        useErrorStore().set(new Error(`${err}`));
      }
    },
    async newNote(key: string) {
      try {
        const note = {
          key,
          body: '',
        }
        const response = await postNote({
          authorization: useTokenStore().authorization,
          requestBody: note,
        })
        if ((response as Message).ok !== undefined) {
          throw new Error((response as Message).message)
        }
        const parts = key.split('/');
        const dir = parts.map((d, i) => {
          if (i === 0 && !d) {
            return '';
          } else if (i < parts.length - 1) {
            return d;
          } else {
            return null;
          }
        }).filter(d => d !== null).join('/') || '/';
        await this.loadDirectory(dir);
        return btoa((response as Note).key);
      } catch (err) {
        useErrorStore().set(new Error(`${err}`));
      }
    },
    async saveNote() {
      try {
        if (!this.note) {
          return
        }
        const response = await putNoteByKey({
          key: btoa(this.note.key),
          authorization: useTokenStore().authorization,
          requestBody: this.note,
        })
        if ((response as Message).ok !== undefined) {
          throw new Error((response as Message).message)
        }
        this.note = (response as Note)
      } catch (err) {
        useErrorStore().set(new Error(`${err}`));
      }
    }
  },
});
