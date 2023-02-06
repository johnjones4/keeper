import React, { useState } from 'react';
import tokenManager from '../models/TokenManager';

interface AddNoteProps {
  onError(err: any): void
  onLogin(): void
}


const Login = (props: AddNoteProps) => {
  const [password, setPassword] = useState('')

  const doLogin = async () => {
    try {
      await tokenManager.login(password)
      props.onLogin()
    } catch (e) {
      props.onError(e)
    }
  }

  return (
    <div>
      <form onSubmit={(e) => {
        e.preventDefault()
        doLogin()
        return false
      }}>
        <h1>Login</h1>
        <div className='InputField'>
          <label htmlFor='password'>Password</label>
          <input className='Input' type='password' value={password} onChange={e => setPassword(e.target.value)} name='password' />
        </div>
        <div className='InputFooter'>
          <button className='Btn AddNoteOk' onClick={() => doLogin()} type='submit'>Login</button>
        </div>
      </form>
    </div>
  )
}

export default Login
