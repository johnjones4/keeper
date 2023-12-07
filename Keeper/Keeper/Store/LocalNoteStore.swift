//
//  LocalNoteStore.swift
//  Keeper
//
//  Created by John Jones on 4/4/23.
//

import Foundation
import SQLite
import SimplePath

class LocalNoteStore {
    var connection: Connection?
    var notes = Table("notes")
    let key = Expression<String>("key")
    let body = Expression<String>("body")
    let modified = Expression<Date>("modified")
    let status = Expression<Status?>("status")
    
    private static var _shared: LocalNoteStore?
    static var shared: LocalNoteStore {
        if let shared = _shared {
            return shared
        }
        let shared = LocalNoteStore()
        _shared = shared
        return shared
    }    

    init() {
        do {
            guard let path = FileManager.default.urls(for: .documentDirectory, in: .userDomainMask).first?.appendingPathComponent("store.sqlite").absoluteString else { return }
            self.connection = try Connection(path)
            try connection?.run(notes.create(ifNotExists: true) { t in
                t.column(key, primaryKey: true)
                t.column(body)
                t.column(modified)
                t.column(status)
            })
        } catch {
            print(error)
            abort()
        }
    }
    
    func add(note: Note) throws {
        guard let connection = self.connection else { return }
        try connection.run(notes.insert(
            key <- note.key,
            body <- note.body,
            modified <- note.modified,
            status <- Status.new
        ))
    }
    
    func update(note: Note) throws {
        guard let connection = self.connection else { return }
        let row = notes.filter(key == note.key)
        try connection.run(row.update(
            body <- note.body,
            modified <- note.modified,
            status <- Status.updated
        ))
    }
    
    func delete(note: Note) throws {
        guard let connection = self.connection else { return }
        let row = notes.filter(key == note.key)
        try connection.run(row.update(
            body <- note.body,
            modified <- note.modified,
            status <- Status.deleted
        ))
    }
    
    func get(key k: String) throws -> Note? {
        guard let connection = self.connection else { return nil }
        let row = notes.filter(key == k)
        for note in try connection.prepare(row) {
            return Note(key: note[key], body: note[body], modified: note[modified], status: note[status])
        }
        return nil
    }
    
    
    func startSyncUp() throws -> [Note] {
        guard let connection = self.connection else { return [] }
        var n : [Note] = []
        for row in try connection.prepare(notes.filter(status != nil)) {
            n.append(Note(key: row[key], body: row[body], modified: row[modified], status: row[status]))
        }
        return n
    }
    
    
    func endSyncUp(items: [Note]) throws {
        //TODO delete
        guard let connection = self.connection else { return }
        for note in items {
            let row = notes.filter(key == note.key)
            try connection.run(row.update(
                status <- nil,
                body <- note.body,
                modified <- note.modified
            ))
        }
    }
    
    func startSyncDown() throws {
        guard let connection = self.connection else { return }
        try connection.run(notes.update(
            status <- Status.syncingDown
        ))
    }
    
    func syncDown(note: Note) throws {
        guard let connection = self.connection else { return }
        let row = notes.filter(key == note.key)
        if try connection.pluck(row) != nil {
            try connection.run(row.update(
                body <- note.body,
                modified <- note.modified,
                status <- nil
            ))
        } else {
            try connection.run(notes.insert(
                key <- note.key,
                body <- note.body,
                modified <- note.modified,
                status <- nil
            ))
        }
    }
    
    func endSyncDown() throws {
        guard let connection = self.connection else { return }
        for item in try connection.prepare(notes.filter(status != nil)) {
            try connection.run(notes.filter(key == item[key]).delete())
        }
    }
    
    func directories(inDir: String) throws -> [Directory] {
        guard let connection = self.connection else { return [] }
        var d = Set<Directory>()
        let depth = Path.split(inDir).count
        for row in try connection.prepare(notes.filter(key.like("\(inDir)%"))) {
            let key = row[key]
            let parts = Path.split(key)
            if parts.count-1 > depth {
                let fullpath = Path.join(Array(parts[0..<depth+1]))
                let dir = parts[depth]
                d.update(with: Directory(fullpath: fullpath, currentPath: dir))
            }
        }
        return d.sorted { a, b in
            return a.currentPath.localizedCompare(b.currentPath) == .orderedAscending
        }
    }
    
    func notes(inDir: String) throws -> [Note] {
        guard let connection = self.connection else { return [] }
        var n = Set<Note>()
        let depth = Path.split(inDir).count
        for row in try connection.prepare(notes.filter(key.like("\(inDir)%"))) {
            let key = row[key]
            let parts = Path.split(key)
            if parts.count == depth+1 {
                n.update(with: Note(key: key, body: row[body], modified: row[modified], status: row[status]))
            }
        }
        return n.sorted { a, b in
            return a.name.localizedCompare(b.name) == .orderedAscending
        }
    }
}
