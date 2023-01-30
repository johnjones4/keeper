//
//  NotesListViewViewModel.swift
//  Keeper
//
//  Created by John Jones on 12/19/22.
//

import Foundation
import Combine

class NoteEditorViewModel: ObservableObject {
    var subscribers = Set<AnyCancellable>()
    var saveTimer: Timer?
    @Published var note: Note
    var repo: NotesManager
    
    init(repo: NotesManager, note: Note) {
        self.repo = repo
        self.note = note
    }
    
    func loadNote() {
        if self.note.id == nil {
            return
        }
        repo
            .getNote(note: note)
            .sink { completion in
                //todo error handling
            } receiveValue: { note in
                self.note = note
            }
            .store(in: &subscribers)
    }
    
    @objc func saveNote() {
        repo
            .saveNote(note: note)
            .sink { completion in
                //todo error handling
            } receiveValue: { note in
                //todo
                print(note)
            }
            .store(in: &subscribers)
    }
    
    func saveNoteLater() {
        if let timer = self.saveTimer {
            timer.invalidate()
        }
        self.saveTimer = Timer.scheduledTimer(timeInterval: 1, target: self, selector: #selector(NoteEditorViewModel.saveNote), userInfo: nil, repeats: false)
    }
}
