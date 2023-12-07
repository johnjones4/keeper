//
//  ContentView.swift
//  Keeper
//
//  Created by John Jones on 4/4/23.
//

import SwiftUI
import Combine

struct NotesListView: View {
    @EnvironmentObject var appViewModel: AppViewModel
    @State var dirs: [Directory] = []
    @State var notes: [Note] = []
    @State var showNewNote: Bool = false
    var dir: String
    
    init(dir: String) {
        self.dir = dir
    }
    
    func refresh(withFetch: Bool) async {
        do {
            if withFetch {
                dirs = try LocalNoteStore.shared.directories(inDir: self.dir)
                notes = try LocalNoteStore.shared.notes(inDir: self.dir)
                try await appViewModel.sync()
            }
        } catch {
            appViewModel.handleError(error: error)
        }
        do {
            dirs = try LocalNoteStore.shared.directories(inDir: self.dir)
            notes = try LocalNoteStore.shared.notes(inDir: self.dir)
        } catch {
            appViewModel.handleError(error: error)
        }
    }
    
    var body: some View {
        List {
            dirs.count == 0 ? nil : Section("Directories") {
                ForEach(dirs) {dir in
                    HStack {
                        Image(systemName: "folder.fill")
                        Button(dir.currentPath) {
                            appViewModel.path.append(dir)
                        }.foregroundColor(Color.label)
                        Spacer()
                        Image(systemName: "chevron.right")
                            .foregroundColor(Color.gray)
                    }
                }
            }
            notes.count == 0 ? nil : Section("Notes") {
                ForEach(notes) {note in
                    HStack {
                        Image(systemName: "note.text")
                        Button(note.name) {
                            appViewModel.path.append(note)
                        }
                        .foregroundColor(Color.label)
                        Spacer()
                        Image(systemName: "chevron.right")
                            .foregroundColor(Color.gray)
                    }
                }
            }
        }.sheet(isPresented: $showNewNote, content: {
            NewNoteView(dir: dir, done: { note in
                appViewModel.path.append(note)
            }).environmentObject(appViewModel)
        })
        .navigationBarTitle(Text(dir))
        .onAppear() {
            Task.init {
                await self.refresh(withFetch: self.dir == "/")
            }
        }
        .refreshable {
            await self.refresh(withFetch: true)
        }
        .toolbar {
            ToolbarItem(placement: .navigationBarTrailing) {
                Button("New") {
                    showNewNote = true
                }
            }
        }
    }
}

struct NotesListView_Previews: PreviewProvider {
    static var previews: some View {
        NotesListView(dir: "/")
    }
}
