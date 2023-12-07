//
//  AppViewModel.swift
//  Keeper
//
//  Created by John Jones on 4/19/23.
//

import Foundation
import Combine
import SwiftUI

class AppViewModel: ObservableObject {
    @Published var path = NavigationPath()
    @Published var isNotReady = !RemoteNoteStore.shared.isReady
    @Published var error: Error?
    @Published var showError = false
    @Published var presentedObjects: [any Hashable] = []
    var cancellable: [AnyCancellable] = []
    var timer: Timer?
    var syncKeys = Set<String>()
    
    func syncSingle(key: String) async throws {
        return try await withCheckedThrowingContinuation { continuation in
            SyncManager.shared.syncSingle(key: key)
                .sink { completion in
                    switch completion {
                    case .failure(let error):
                        continuation.resume(throwing: error)
                    case .finished:
                        continuation.resume(returning: ())
                    }
                } receiveValue: { _ in
                    
                }
                .store(in: &self.cancellable)
        }
    }
    
    func sync() async throws {
        return await withCheckedContinuation({ continuation in
            SyncManager.shared.syncUp()
                .sink { completion in
                    switch completion {
                    case .failure(let error):
                        self.handleError(error: error)
                    case .finished:
                        break
                    }
                    SyncManager.shared.syncDown()
                        .sink { completion in
                            switch completion {
                            case .failure(let error):
                                self.handleError(error: error)
                            case .finished:
                                continuation.resume(returning: ())
                            }
                            
                        } receiveValue: { _ in
                            
                        }
                        .store(in: &self.cancellable)
                } receiveValue: { _ in
                    
                }
                .store(in: &self.cancellable)
        })
//        return try await withCheckedThrowingContinuation { continuation in
            
//        }
    }
    
    private func _enqueueSingleSync() {
        if let timer = timer {
            timer.invalidate()
        }
        self.timer = Timer.scheduledTimer(withTimeInterval: 1, repeats: false) { t in
            self.timer = nil
            guard let key = self.syncKeys.popFirst() else {
                return
            }
            SyncManager.shared.syncSingle(key:key)
                .sink { completion in
                    switch completion {
                    case .failure(let error):
                        self.handleError(error: error)
                    case .finished:
                        break
                    }
                    self._enqueueSingleSync()
                } receiveValue: { _ in
                    
                }
                .store(in: &self.cancellable)
        }
    }
    
    func enqueueSingleSync(key: String) {
        syncKeys.update(with: key)
        _enqueueSingleSync()
    }
    
    func updateReadyState() {
        self.isNotReady = !RemoteNoteStore.shared.isReady
    }
    
    func handleError(error: Error) {
        if let rerror = error as? RemoteError {
            switch rerror {
            case .AccessDenied:
                CredentialStore.shared.apiToken = nil
                updateReadyState()
            default:
                self.error = error
                self.showError = true
                break
            }
        }
    }
    
    public func pushView(_ destination: any Hashable) {
        self.path.append(destination)
    }
}
