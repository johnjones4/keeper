//
//  Note.swift
//  Keeper
//
//  Created by John Jones on 12/19/22.
//

import Foundation

struct NoteBody: Codable, Hashable, Equatable {
    var text: String = ""
    let structuredData: StructuredDataProperty?
}

struct Note: Codable, Identifiable, Hashable, Equatable {
    let id: String?
    var path: String?
    var title: String = ""
    var body: NoteBody = NoteBody()
    var tags: [String] = []
    var sourceURL: String?
    var source: String = "app"
    var format: String?
//    var created: Date?
//    var updated: Date?
    
    init(id: String, title: String, bodyText: String) {
        self.id = id
        self.title = title
        self.body = NoteBody(text: bodyText)
    }
    
    init() {
        self.id = nil
        self.title = "New note"
        self.format = "text/markdown"
        self.body = NoteBody(text: "")
    }
}

struct NotesResponse: Decodable {
    let items: [Note]
}
