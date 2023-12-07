//
//  RemoteNoteStore.swift
//  Keeper
//
//  Created by John Jones on 4/7/23.
//

import Foundation
import Combine

struct LoginRequest: Encodable {
    let password: String
}

struct LoginResponse: Decodable {
    let token: String
}

struct ErrorResponse: Decodable {
    let message: String
}

enum RemoteError: Error {
    case AccessDenied
    case Unknown(status: Int, message: String)
}

class RemoteNoteStore {
    let apiRoot: String = "https://notes.johnjonesfour.com"
//    let apiRoot: String = "http://localhost:8080"
    var authKey: String? {
        return CredentialStore.shared.apiToken
    }
    
    var isReady: Bool {
        return authKey != nil
    }
    
    private static var _shared: RemoteNoteStore?
    static var shared: RemoteNoteStore {
        if let shared = _shared {
            return shared
        }
        let shared = RemoteNoteStore()
        _shared = shared
        return shared
    }
    
    private static var dateFormatter: DateFormatter {
        let dateFormatter = DateFormatter()
        dateFormatter.dateFormat = "yyyy-MM-dd'T'HH:mm:ssxxx"
        dateFormatter.locale = Locale(identifier: "en_US")
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
                request.httpBody = try RemoteNoteStore.encoder.encode(b)
                request.setValue("application/json", forHTTPHeaderField: "Content-type")
            }
            
            if let authKey = self.authKey {
                request.setValue("Bearer \(authKey)", forHTTPHeaderField: "Authorization")
            }
            
            return URLSession.shared.dataTaskPublisher(for: request)
                .tryMap { data, response in
                    if let httpResponse = response as? HTTPURLResponse {
                        if httpResponse.statusCode == 401 {
                            throw RemoteError.AccessDenied
                        } else if httpResponse.statusCode >= 300 {
                            throw RemoteError.Unknown(status: httpResponse.statusCode, message: "")
                        }
                    }
                    return data

                }
                .mapError { $0 as Error }
                .decode(type: V.self, decoder: RemoteNoteStore.decoder)
                .eraseToAnyPublisher()
                .receive(on: RunLoop.main)
                .eraseToAnyPublisher()
        } catch {
            return Fail<V, Error>(error: error)
                .eraseToAnyPublisher()
        }
    }
    
    func login(password: String) -> AnyPublisher<LoginResponse, Error> {
        return request(method: "POST", path: "/api/token", body: LoginRequest(password: password))
    }
    
    func saveNote(note: Note, isNew: Bool) -> AnyPublisher<Note, Error> {
        guard let safeId = note.id.addingPercentEncoding(withAllowedCharacters: CharacterSet.urlPathAllowed) else {
            return Fail<Note, Error>(error: NSError(domain: "", code: 0)).eraseToAnyPublisher()
        }
        let path = isNew ? "/api/note" : "/api/note/\(safeId)"
        let method = isNew ? "POST" : "PUT"
        return request(method: method, path: path, body: note)
    }
    
    func getNotes(page: String) -> AnyPublisher<NotesResponse, Error> {
        return request(method: "GET", path: "/api/note?page=\(page)", body: nil as String?)
    }
    
    func getNote(key: String) -> AnyPublisher<Note, Error> {
        let id = Data(key.utf8).base64EncodedString()
        guard let safeId = id.addingPercentEncoding(withAllowedCharacters: CharacterSet.urlPathAllowed) else {
            return Fail<Note, Error>(error: NSError(domain: "", code: 0)).eraseToAnyPublisher()
        }
        let path = "/api/note/\(safeId)"
        return request(method: "GET", path: path, body: nil as String?)
    }
    
    func getAllNotes() -> AnyPublisher<[String], Error> {
        return getNextNotes(page: "", notes: [])
    }
    
    private func getNextNotes(page: String, notes: [String]) -> AnyPublisher<[String], Error> {
        return getNotes(page: page)
            .flatMap { response -> AnyPublisher<[String], Error> in
                if response.nextPage != "" {
                    return self.getNextNotes(page: response.nextPage, notes: notes + response.notes).eraseToAnyPublisher()
                } else {
                    return Just(notes + response.notes)
                        .setFailureType(to: Error.self)
                        .eraseToAnyPublisher()
                }
            }
            .eraseToAnyPublisher()
    }
}

struct NotesResponse: Decodable {
    let notes: [String]
    let nextPage: String
}
