//
//  NewNoteView.swift
//  Keeper
//
//  Created by John Jones on 4/20/23.
//

import SwiftUI

struct NewNoteView: View {
    typealias Done = ((Note) -> Void)
    @EnvironmentObject var appViewModel: AppViewModel
    @ObservedObject var viewModel: NewNoteViewModel
    @Environment(\.dismiss) private var dismiss
    var done: Done
    
    
    init(dir: String, done: @escaping Done) {
        self.viewModel = NewNoteViewModel(dir: dir)
        self.done = done
    }
    
    var body: some View {
        VStack {
            Text("New Note")
                .font(.largeTitle)
                .fontWeight(.bold)
                .padding(.bottom, 20)
            LabeledContent {
                TextField("Directory", text: $viewModel.dir, prompt: Text("Directory"))
                    .textFieldStyle(RoundedBorderTextFieldStyle())
            } label: {
                Text("Directory")
            }
            LabeledContent {
                TextField("Name", text: $viewModel.name, prompt: Text("Name"))
                    .textFieldStyle(RoundedBorderTextFieldStyle())
            } label: {
                Text("Name")
            }
            LabeledContent {
                Picker("Type", selection: $viewModel.ext) {
                    ForEach(NoteType.allCases) { type in
                        Text(type.rawValue).tag(type)
                    }
                }
            } label: {
                Text("Type")
            }
            
            HStack {
                Button(action: {
                    dismiss()
                }) {
                    Text("Cancel")
                        .font(.subheadline)
                        .foregroundColor(.white)
                        .padding()
                        .background(Color.gray)
                        .cornerRadius(10)
                }
                Button(action: {
                    Task.init {
                        do {
                            let key = try viewModel.createNote()
                            try await appViewModel.syncSingle(key: key)
                            if let note = try LocalNoteStore.shared.get(key: key) {
                                self.done(note)
                            }
                        } catch {
                            appViewModel.handleError(error: error)
                        }
                        dismiss()
                    }
                }) {
                    Text("Create")
                        .font(.subheadline)
                        .foregroundColor(.white)
                        .padding()
                        .background(Color.blue)
                        .cornerRadius(10)
                }
            }
        }.padding()
    }
}

struct NewNoteView_Previews: PreviewProvider {
    static var previews: some View {
        NewNoteView(dir: "/") { _ in
            
        }
    }
}
