//
//  NoteEditor.swift
//  Keeper
//
//  Created by John Jones on 12/19/22.
//

import SwiftUI
import CodeEditor

struct NoteTextEditor: View {
    @StateObject var viewModel: NoteEditorViewModel
    @State var fontSize: CGFloat = 18
    @Environment(\.colorScheme) var colorScheme
    var theme : CodeEditor.ThemeName {
        return self.colorScheme == .light ? CodeEditor.ThemeName(rawValue: "xcode") : CodeEditor.ThemeName(rawValue: "defaultt")
    }
    
    var body: some View {
        VStack {
            TextField("Title", text: $viewModel.note.title)
                .padding(EdgeInsets(top: 10, leading: 10, bottom: 10, trailing: 10))
                .font(Font.system(size: 24))
            CodeEditor(source: $viewModel.note.body.text, language: .markdown, theme: .default, fontSize: $fontSize, inset: CGSize(width: 10, height: 10))
        }
            .onChange(of: viewModel.note.body.text, perform: { newValue in
                viewModel.saveNoteLater()
            })
            .onChange(of: viewModel.note.title, perform: { newValue in
                viewModel.saveNoteLater()
            })
            .onDisappear() {
                viewModel.saveNote()
            }
            .onAppear() {
                viewModel.loadNote()
            }
#if os(iOS)
            .navigationBarTitleDisplayMode(NavigationBarItem.TitleDisplayMode.inline)
#endif
    }
}

struct NoteTextEditor_Previews: PreviewProvider {
    static var previews: some View {
        NoteTextEditor(viewModel: NoteEditorViewModel(repo: TestNotes(), note: Note(id: "test", title: "title", bodyText: "text")))
    }
}
