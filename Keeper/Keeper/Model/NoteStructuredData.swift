//
//  NoteStructuredData.swift
//  Keeper
//
//  Created by John Jones on 12/22/22.
//

import Foundation

struct NoteSection {
    let title: String
    let rows: [String]
}

extension Note {
    var sections: [NoteSection] {
        return []
    }
    
    var searchProperties: [String] {
        
    }
}

