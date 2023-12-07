//
//  LoginView.swift
//  Keeper
//
//  Created by John Jones on 4/19/23.
//

import SwiftUI

struct LoginView: View {
    @EnvironmentObject var appViewModel: AppViewModel
    @ObservedObject var viewModel = LoginViewModel()
    
    var body: some View {
        VStack {
            Text("Login")
                .font(.largeTitle)
                .fontWeight(.bold)
                .padding(.bottom, 50)
            
            Image(systemName: "person.circle.fill")
                .font(.system(size: 100))
                .foregroundColor(.gray)
                .padding(.bottom, 50)
            
            SecureField("Password", text: $viewModel.password)
                .textFieldStyle(RoundedBorderTextFieldStyle())
                .padding(.horizontal, 50)
            
            Button(action: {
                viewModel.login { error in
                    if let error = error {
                        appViewModel.handleError(error: error)
                    } else {
                        appViewModel.updateReadyState()
                    }
                }
            }) {
                Text("Login")
                    .font(.headline)
                    .foregroundColor(.white)
                    .padding()
                    .frame(width: 200, height: 50)
                    .background(Color.blue)
                    .cornerRadius(10)
            }
        }.padding()
    }
}

struct LoginView_Previews: PreviewProvider {
    static var previews: some View {
        LoginView()
    }
}
