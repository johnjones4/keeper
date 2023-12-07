//
//  LoginViewModel.swift
//  Keeper
//
//  Created by John Jones on 4/19/23.
//

import Foundation
import Combine

class LoginViewModel: ObservableObject {
    @Published var password: String = ""
    var cancellable: [AnyCancellable] = []
    
    func login(complete: @escaping ((Error?) -> Void)) {
        RemoteNoteStore.shared.login(password: password)
            .sink { completion in
                switch completion {
                case .failure(let error):
                    complete(error)
                case .finished:
                    complete(nil)
                }
            } receiveValue: { response in
                CredentialStore.shared.apiToken = response.token
            }
            .store(in: &cancellable)
    }
}
