//
//  NewNoteViewModel.swift
//  Keeper
//
//  Created by John Jones on 4/20/23.
//

import Foundation
import Combine
import SimplePath

enum NoteType: String, Identifiable {
    case Text = "Text"
    case Markdown = "Markdown"
    case Todo = "Todo"
    
    static var allCases: [NoteType] {
        return [.Text, .Markdown, .Todo]
    }
    
    var id: String {
        return self.rawValue
    }
    
    var ext: String {
        switch self {
        case .Text:
            return "txt"
        case .Markdown:
            return "md"
        case .Todo:
            return "todo"
        }
    }
}

class NewNoteViewModel: ObservableObject {
    @Published var dir: String
    @Published var name: String
    @Published var ext: NoteType = .Text
    
    private static var dateFormatter: DateFormatter {
        let dateFormatter = DateFormatter()
        dateFormatter.dateFormat = "yyyy-MM-dd'T'HH:mm:ssxxx"
        dateFormatter.locale = Locale(identifier: "en_US")
        return dateFormatter
    }
    
    init(dir: String) {
        self.dir = dir
        self.name = "new note \(NewNoteViewModel.dateFormatter.string(from: Date()))"
    }
    
    func createNote() throws -> String {
        let key = Path.join([dir, name + "." + ext.ext])
        try LocalNoteStore.shared.add(note: Note(key: key, body: "New Note", modified: Date(), status: nil))
        return key
    }
}
