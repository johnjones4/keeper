//
//  StructuredDataProperty.swift
//  Keeper
//
//  Created by John Jones on 12/22/22.
//

import Foundation

struct StructuredDataProperty: Codable {
    let type: [String]
    let str: String
    let int: Int
    let float: Double
    let bool: Bool
    let id: String
    let properties: [StructuredDataProperty]
}
