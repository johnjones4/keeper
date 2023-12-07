//
//  Note.swift
//  Keeper
//
//  Created by John Jones on 4/4/23.
//

import Foundation
import SQLite
import SimplePath

enum Status: String, Codable, Value {
    static func fromDatatypeValue(_ datatypeValue: String) -> Status {
        return Status(rawValue: datatypeValue)!
    }
    
    var datatypeValue: String {
        return self.rawValue
    }
    
    typealias Datatype = String
    static var declaredDatatype = "TEXT"
    
    case new = "new"
    case updated = "updated"
    case deleted = "deleted"
    case syncingDown = "syncingdown"
}

struct Directory: Identifiable, Hashable {
    let fullpath: String
    let currentPath: String
    var id: String {
        return fullpath
    }
}

struct Note: Codable, Identifiable, Hashable {
    let key: String
    let body: String
    let modified: Date
    let status: Status?
    
    var id: String {
        return Data(key.utf8).base64EncodedString()
    }
    
    var name: String {
        return Path.basename(key)
    }
}
