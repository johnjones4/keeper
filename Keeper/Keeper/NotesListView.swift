//
//  ContentView.swift
//  Shared
//
//  Created by John Jones on 8/15/22.
//

import SwiftUI
import Combine

class NotesListViewViewModel: ObservableObject {
    @Published var notes: [Note] = []
    @Published var error: Error?
    private var subscribers = Set<AnyCancellable>()
    
    let repo: NotesManager
    
    init(repo: NotesManager) {
        self.repo = repo
    }
    
    func deleteNoteAt(index: Int) {
        repo.deleteNote(note: notes[index])
            .sink(
                receiveCompletion: { (completion) in
                    switch completion {
                    case .finished:
                        break
                    case .failure(let error):
                        self.error = error
                    }
                },
                receiveValue: {_ in
                    self.loadNotes()
                }
            )
            .store(in: &subscribers)
    }
    
    func loadNotes() {
        repo.getNotes()
            .sink(
                receiveCompletion: { (completion) in
                    switch completion {
                    case .finished:
                        break
                    case .failure(let error):
                        self.error = error
                    }
                },
                receiveValue: {
                    self.notes = $0.items
                }
            )
            .store(in: &subscribers)

    }
}

struct NotesListView: View {
    @ObservedObject var viewModel: NotesListViewViewModel
    var factory = EditorFactory()
    
    var body: some View {
        NavigationView {
            List {
                ForEach(viewModel.notes) {note in
                    NavigationLink(destination: factory.editor(viewModel: NoteEditorViewModel(repo: self.viewModel.repo, note: note))) {
                        VStack(alignment: .leading) {
                            Text(note.title)
                        }
                    }
                }
                .onDelete { indexSet in
                    indexSet.forEach { i in
                        viewModel.deleteNoteAt(index: i)
                    }
                }
            }
            .navigationTitle("Notes")
            .toolbar {
                ToolbarItem(placement: .primaryAction) {
                    NavigationLink(destination: factory.editor(viewModel: NoteEditorViewModel(repo: self.viewModel.repo, note: Note()))) {
                        Text("New")
                    }

                }
            }
            .onAppear {
                    self.viewModel.loadNotes()
            }
        }
    }
}

struct NotesListView_Previews: PreviewProvider {
    static var previews: some View {
        NotesListView(
            viewModel: NotesListViewViewModel(
                repo: TestNotes()
            )
        )
    }
}
