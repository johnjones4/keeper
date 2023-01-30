//
//  EditorFactory.swift
//  Keeper
//
//  Created by John Jones on 12/19/22.
//

import Foundation
import SwiftUI

class EditorFactory {
    func editor(viewModel: NoteEditorViewModel) -> (some View)? {
        if (viewModel.note.body.text != "") {
            return NoteTextEditor(viewModel: viewModel)
        }
        return NoteTextEditor(viewModel: viewModel) //TODO
    }
}
