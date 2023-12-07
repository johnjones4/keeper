//
//  NoteView.swift
//  Keeper
//
//  Created by John Jones on 4/18/23.
//

import SwiftUI
import Combine
import SimplePath

struct NoteView: View {
    @EnvironmentObject private var appViewModel: AppViewModel
    @StateObject private var viewModel: NoteViewModel
    
    init(note: Note) {
        _viewModel = StateObject(wrappedValue: {NoteViewModel(note: note)}())
    }
    
    var editor: some View {
        switch Path.extname(viewModel.key) {
        case "todo":
            return AnyView(TodoEditorView(text: $viewModel.body))
        default:
            return AnyView(TextEditor(text: $viewModel.body))
        }
    }
    
    var body: some View {
        editor
            .onChange(of: viewModel.body, perform: { newValue in
                viewModel.waitAndChange { error in
                    if let error = error {
                        self.appViewModel.handleError(error: error)
                    } else {
                        appViewModel.enqueueSingleSync(key: viewModel.key)
                    }
                }
            })
            .onAppear() {
                
            }
            .navigationBarTitle(Text(viewModel.key))
    }
}

struct NoteView_Previews: PreviewProvider {
    static var previews: some View {
        NoteView(note: Note(key: "/test.text", body: "dfsdfsdfsdf", modified: Date(), status: nil))
    }
}
