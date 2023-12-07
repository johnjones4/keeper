//
//  NoteViewModel.swift
//  Keeper
//
//  Created by John Jones on 4/19/23.
//

import Foundation
import Combine

class NoteViewModel: ObservableObject {
    @Published var key: String
    @Published var body: String
    @Published var modified: Date
    var cancellable: [AnyCancellable] = []
    var timer: Timer?
    
    init(note: Note) {
        self.key = note.key
        self.body = note.body
        self.modified = note.modified
    }
    
    func waitAndChange(complete: @escaping ((Error?) -> Void)) {
        if let timer = timer {
            timer.invalidate()
        }
        self.timer = Timer.scheduledTimer(withTimeInterval: 0.5, repeats: false) { t in
            do {
                try LocalNoteStore.shared.update(note: Note(key: self.key, body: self.body, modified: self.modified, status: nil))
                complete(nil)
            } catch {
                complete(error)
            }
        }
    }
}
