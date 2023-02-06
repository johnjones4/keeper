interface TokenResponse {
  token: string
}

class TokenManager {
  key: string
  token: string

  constructor(key: string) {
    this.key = key
    const t = localStorage.getItem(this.key) 
    this.token = t ? t : ''
  }

  getToken(): string {
    return this.token
  }

  isReady() {
    return this.getToken() !== ''
  }

  async login(password: string) {
    const res = await fetch('/api/token', {
      method: 'POST',
      body: JSON.stringify({password})
    })
    const json = await res.json()
    if (res.status !== 200) {
      throw new Error(json.message)
    }
    const tokenInfo = json as TokenResponse
    localStorage.setItem(this.key, tokenInfo.token)
    this.token = tokenInfo.token
  }

  getHeaders(): any {
    return {
      'Authorization': `Bearer ${this.getToken()}`
    }
  }

  clear() {
    this.token = ''
    localStorage.setItem(this.key, '')
  }
}

const tokenManager = new TokenManager('token')
export default tokenManager
