//
//  TodoEditorView.swift
//  Keeper
//
//  Created by John Jones on 4/21/23.
//

import SwiftUI

struct Todo: Identifiable {
    var id: String {
        return label
    }
    
    let complete: Bool
    let label: String
    let indent: Int
    
    init(indent: Int, complete: Bool, label: String) {
        self.indent = indent
        self.complete = complete
        self.label = label
    }
    
    init(line: String) {
        var spaces = 0
        for char in line {
            if char == " " {
                spaces+=1
            } else {
                break
            }
        }
        self.indent = spaces
        self.complete = line.substring(from: line.index(line.startIndex, offsetBy: spaces)).hasPrefix("[x] ")
        self.label = line.substring(from: line.index(line.startIndex, offsetBy: spaces+4))
    }
    
    var toggled: Todo {
        return Todo(indent: indent, complete: !complete, label: label)
    }
    
    func indented(_ direction: Int) -> Todo {
        return Todo(indent: max(0, indent + direction), complete: complete, label: label)
    }
 
    var string: String {
        var str = ""
        var i = 0
        while i < self.indent {
            str += " "
            i+=1
        }
        str += self.complete ? "[x] " : "[ ] "
        str += self.label
        return str
    }
}

struct TodoEditorView: View {
    @Binding var text: String
    @State var newTodo: String = ""
    
    var todos: [Todo] {
        return self.text.split(separator: "\n").enumerated().map { (n, str) -> Todo in
            return Todo(line: String(str))
        }
    }
    
    func saveTodos(_ todos: [Todo]) {
        text = todos.map({ todo in todo.string }).joined(separator: "\n")
    }
    
    func submitNewTodo() {
        saveTodos(todos + [Todo(indent: 0, complete: false, label: newTodo)])
        newTodo = ""
    }
    
    var body: some View {
        List {
            ForEach(todos) { todo in
                HStack {
                    Button {
                        self.saveTodos(todos.map({ todo1 in
                            if todo.label == todo1.label {
                                return todo1.toggled
                            }
                            return todo1
                        }))
                    } label: {
                        HStack {
                            Image(systemName: todo.complete ? "checkmark.circle.fill" : "circle")
                            Text(todo.label).strikethrough(todo.complete)
                                .foregroundColor(Color.label)
                        }.padding(EdgeInsets(top: 0, leading: CGFloat(todo.indent) * 10, bottom: 0, trailing: 0))
                    }
                    Spacer()
                    Image(systemName: "arrowshape.left.fill")
                        .onTapGesture {
                            if todo.indent > 0 {
                                self.saveTodos(todos.map({ todo1 in
                                    if todo.label == todo1.label {
                                        return todo1.indented(-1)
                                    }
                                    return todo1
                                }))
                            }
                        }
                        .foregroundColor(todo.indent == 0 ? Color.gray : Color.primary)
                    Image(systemName: "arrowshape.right.fill")
                        .onTapGesture {
                            self.saveTodos(todos.map({ todo1 in
                                if todo.label == todo1.label {
                                    return todo1.indented(1)
                                }
                                return todo1
                            }))
                        }
                        .foregroundColor(Color.primary)
                }
            }
            .onDelete { idx in
                self.saveTodos(todos.enumerated().filter({ (index, _) in
                    return !idx.contains(index)
                }).map({ (_, item) in item }))
            }
            HStack {
                TextField("New Todo", text: $newTodo)
                    .onSubmit {
                        submitNewTodo()
                    }
                Button {
                    submitNewTodo()
                } label: {
                    Image(systemName: "plus.app.fill").foregroundColor(Color.green)
                }

            }
        }
    }
}

//struct TodoEditorView_Previews: PreviewProvider {
//    static var previews: some View {
//        TodoEditorView(text: "")
//    }
//}
