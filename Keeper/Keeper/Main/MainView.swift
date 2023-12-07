//
//  MainView.swift
//  Keeper
//
//  Created by John Jones on 4/19/23.
//

import SwiftUI

struct MainView: View {
    @StateObject private var appViewModel = AppViewModel()
    
    var body: some View {
        ZStack {
            NavigationStack(path: $appViewModel.path) {
                NotesListView(dir: "/").environmentObject(appViewModel)
                    .environmentObject(appViewModel)
                    .navigationDestination(for: Note.self, destination: { note in
                        NoteView(note: note)
                            .environmentObject(appViewModel)
                            .environmentObject(appViewModel)
                    })
                    .navigationDestination(for: Directory.self, destination: { directory in
                        NotesListView(dir: directory.fullpath)
                            .environmentObject(appViewModel)
                            .environmentObject(appViewModel)
                    })
            }
            .alert("Error", isPresented: $appViewModel.showError, actions: {
                
            }, message: {
                Text(appViewModel.error?.localizedDescription ?? "")
            })
            .fullScreenCover(isPresented: $appViewModel.isNotReady) {
                Task.init {
                    do {
                        try await appViewModel.sync()
                    } catch {
                        appViewModel.handleError(error: error)
                    }
                }
            } content: {
                LoginView()
                    .environmentObject(appViewModel)
            }
        }
        .toolbar {
            !appViewModel.isNotReady ? ToolbarItem(placement: .navigationBarLeading) {
                Button("Logout") {
                    CredentialStore.shared.apiToken = nil
                    appViewModel.updateReadyState()
                }
            } : nil
        }
    }
}

struct MainView_Previews: PreviewProvider {
    static var previews: some View {
        MainView()
    }
}
