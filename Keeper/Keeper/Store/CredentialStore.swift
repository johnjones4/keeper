//
//  CredentialStore.swift
//  Keeper
//
//  Created by John Jones on 4/19/23.
//

import Foundation

class CredentialStore {
    private static let service = "com.johnjonesfour.keeper"
    private static var _shared: CredentialStore?
    static var shared: CredentialStore {
        if let shared = _shared {
            return shared
        }
        let shared = CredentialStore()
        _shared = shared
        return shared
    }
    
    var apiToken: String? {
        get {
            return getValue(key: "apitoken")
        }
        set {
            if let value = newValue {
                setValue(key: "apitoken", value: value)
            } else {
                removeValue(key: "apitoken")
            }
        }
    }
    
    private func setValue(key: String, value: String) {
        removeValue(key: key)
        let addquery: [String: Any] = [kSecClass as String: kSecClassGenericPassword,
                                       kSecAttrService as String: CredentialStore.service as AnyObject,
                                       kSecAttrAccount as String: key as AnyObject,
                                       kSecValueData as String: value.data(using: .utf8) as AnyObject]
        
        SecItemAdd(addquery as CFDictionary, nil)
    }
    
    private func getValue(key: String) -> String? {
        let getquery: [String: Any] = [kSecClass as String: kSecClassGenericPassword,
                                       kSecAttrService as String: CredentialStore.service as AnyObject,
                                       kSecAttrAccount as String: key as AnyObject,
                                       kSecMatchLimit as String: kSecMatchLimitOne,
                                       kSecReturnData as String: kCFBooleanTrue as AnyObject]
        var item: AnyObject?
        let status = SecItemCopyMatching(getquery as CFDictionary, &item)
        guard status == errSecSuccess else { return nil }
        guard let data = item as? Data else {
            return nil
        }
        return String(data: data, encoding: .utf8)
    }
    
    private func removeValue(key: String) {
        let query: [String: Any] = [kSecClass as String: kSecClassGenericPassword,
                                    kSecAttrService as String: CredentialStore.service as AnyObject,
                                    kSecAttrAccount as String: key as AnyObject]
        SecItemDelete(query as CFDictionary)
    }
}
