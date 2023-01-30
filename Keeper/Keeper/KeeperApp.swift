//
//  KeeperApp.swift
//  Keeper
//
//  Created by John Jones on 12/19/22.
//

import SwiftUI

#if DEBUG
let root = "http://localhost:8080"
#else
let root = "https://keeper.johnjonesfour.com"
#endif

@main
struct KeeperApp: App {
    var body: some Scene {
        WindowGroup {
            NotesListView(
                viewModel: NotesListViewViewModel(
                    repo: HttpNotesManager(apiRoot: root)
                )
            )
        }
    }
}
