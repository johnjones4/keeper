//
//  NotesManager.swift
//  Keeper
//
//  Created by John Jones on 12/19/22.
//

import Foundation
import Combine

protocol NotesManager {
    func saveNote(note: Note) -> AnyPublisher<Note, Error>
    func getNotes() -> AnyPublisher<NotesResponse, Error>
    func getNote(note: Note) -> AnyPublisher<Note, Error>
    func deleteNote(note: Note) -> AnyPublisher<Note, Error>
}

struct HttpNotesManager: NotesManager {
    let apiRoot: String
    
    private static var dateFormatter: DateFormatter {
        let dateFormatter = DateFormatter()
        dateFormatter.dateFormat = "yyyy-MM-dd'T'HH:mm:ss.SSSxxx"
        dateFormatter.locale = Locale(identifier: "en_US")
        //        dateFormatter.timeZone = TimeZone(secondsFromGMT: 0)
        return dateFormatter
    }
    
    private static var decoder: JSONDecoder {
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .formatted(dateFormatter)
        return decoder
    }
    
    private static var encoder: JSONEncoder {
        let encoder = JSONEncoder()
        encoder.dateEncodingStrategy = .formatted(dateFormatter)
        return encoder
    }
    
    private func request<T, V>(method: String, path: String, body: T?) -> AnyPublisher<V, Error> where T : Encodable, V : Decodable {
        do {
            print("(\(method)) \(path)")
            var request = URLRequest(url: URL(string: self.apiRoot+path)!)
            request.httpMethod = method
            
            if let b = body {
                request.httpBody = try HttpNotesManager.encoder.encode(b)
                request.setValue("application/json", forHTTPHeaderField: "Content-type")
            }
            
            return URLSession.shared.dataTaskPublisher(for: request)
                .map { (data: Data, response: URLResponse) -> Data in
                    print(String(data: data, encoding: .utf8) ?? "")
                    return data
                }
                .decode(type: V.self, decoder: HttpNotesManager.decoder)
                .eraseToAnyPublisher()
                .receive(on: RunLoop.main)
                .eraseToAnyPublisher()
        } catch {
            return Fail<V, Error>(error: error)
                .eraseToAnyPublisher()
        }
    }
    
    private func emptyRequest(method: String, path: String) -> AnyPublisher<Void, Error> {
        print("(\(method)) \(path)")
        var request = URLRequest(url: URL(string: self.apiRoot+path)!)
        request.httpMethod = method
        
        return URLSession.shared.dataTaskPublisher(for: request)
            .map { (data: Data, response: URLResponse) in
                print(String(data: data, encoding: .utf8) ?? "")
                return ()
            }
            .mapError({ f in
                return f
            })
            .eraseToAnyPublisher()
            .receive(on: RunLoop.main)
            .eraseToAnyPublisher()
    }
    
    func saveNote(note: Note) -> AnyPublisher<Note, Error> {
        let path = "/api/note" + (note.id == nil ? "" : "/\(note.id!)")
        let method = note.id == nil ? "POST" : "PUT"
        return request(method: method, path: path, body: note)
    }
    
    func getNotes() -> AnyPublisher<NotesResponse, Error> {
        return request(method: "GET", path: "/api/note", body: nil as String?)
    }
    
    func getNote(note: Note) -> AnyPublisher<Note, Error> {
        return request(method: "GET", path: "/api/note/\(note.id!)", body: nil as String?)
    }
    
    func deleteNote(note: Note) -> AnyPublisher<Note, Error> {
        return request(method: "DELETE", path: "/api/note/\(note.id!)", body: nil as String?)
    }
}

struct TestNotes: NotesManager {
    func saveNote(note: Note) -> AnyPublisher<Note, Error> {
        return Just(note)
            .setFailureType(to: Error.self)
            .eraseToAnyPublisher()
    }
    
    func getNote(note: Note) -> AnyPublisher<Note, Error> {
        return Just(note)
            .setFailureType(to: Error.self)
            .eraseToAnyPublisher()
    }
    
    func getNotes() -> AnyPublisher<NotesResponse, Error> {
        return Just(NotesResponse(items: [
            Note(id: "test1", title: "Test 1", bodyText: "test test test 1"),
            Note(id: "test2", title: "Test 2", bodyText: "test test test 2")
        ]))
            .setFailureType(to: Error.self)
            .eraseToAnyPublisher()
    }
    
    func deleteNote(note: Note) -> AnyPublisher<Note, Error> {
        return Just(note)
            .setFailureType(to: Error.self)
            .eraseToAnyPublisher()
    }
    
    
}
