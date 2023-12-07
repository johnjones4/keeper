//
//  SyncManager.swift
//  Keeper
//
//  Created by John Jones on 4/7/23.
//

import Foundation
import Combine

class SyncManager {
    private let remoteStore: RemoteNoteStore
    private let localStore: LocalNoteStore
    
    private static var _shared: SyncManager?
    static var shared: SyncManager {
        if let shared = _shared {
            return shared
        }
        let shared = SyncManager(remoteStore: RemoteNoteStore.shared, localStore: LocalNoteStore.shared)
        _shared = shared
        return shared
    }
    
    init(remoteStore: RemoteNoteStore, localStore: LocalNoteStore) {
        self.remoteStore = remoteStore
        self.localStore = localStore
    }
    
    func syncSingle(key: String) -> AnyPublisher<Void, Error> {
        do {
            guard let note = try self.localStore.get(key: key) else {
                return Fail<Void, Error>(error: NSError(domain: "", code: 0)).eraseToAnyPublisher()
            }
            return self.syncUpStep(notes: [note])
        } catch {
            return Fail<Void, Error>(error: error)
                .eraseToAnyPublisher()
        }
    }
    
    func syncUp() -> AnyPublisher<Void, Error> {
        do {
            let notes = try self.localStore.startSyncUp()
            return self.syncUpStep(notes: notes)
        } catch {
            return Fail<Void, Error>(error: error)
                .eraseToAnyPublisher()
        }
    }
    
    private func syncUpStep(notes: [Note]) -> AnyPublisher<Void, Error> {
        var _notes = notes;
        guard let note = _notes.popLast() else {
            return Just<Void>(()).setFailureType(to: Error.self).eraseToAnyPublisher()
        }
        switch note.status {
        case .new, .updated:
            return remoteStore.saveNote(note: note, isNew: note.status == .new)
                .flatMap { _ -> AnyPublisher<Void, Error> in
                    return self.syncUpStep(notes: _notes)
                }
                .eraseToAnyPublisher()
        default:
            return self.syncUpStep(notes: _notes)
        }
    }
    
    func syncDown() -> AnyPublisher<Void, Error> {
        do {
            try self.localStore.startSyncDown()
        } catch {
            return Fail(error: error).eraseToAnyPublisher()
        }
        return remoteStore.getAllNotes()
            .flatMap { noteDirs -> AnyPublisher<Void, Error> in
                return self.syncDownStep(paths: noteDirs)
            }
            .eraseToAnyPublisher()

    }
    
    private func syncDownStep(paths: [String]) -> AnyPublisher<Void, Error> {
        var _paths = paths
        guard let path = _paths.popLast() else {
            do {
                try self.localStore.endSyncDown()
                return Just<Void>(()).setFailureType(to: Error.self).eraseToAnyPublisher()
            } catch {
                return Fail(error: error).eraseToAnyPublisher()
            }
        }
        return self.remoteStore.getNote(key: path)
            .flatMap { note -> AnyPublisher<Void, Error> in
                do {
                    try self.localStore.syncDown(note: note)
                    return self.syncDownStep(paths: _paths)
                } catch {
                    return Fail(error: error).eraseToAnyPublisher()
                }
            }
            .eraseToAnyPublisher()
    }
}
