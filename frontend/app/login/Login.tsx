// components/Login.js
import React from 'react';
import './Login.css'

export default function Login() {
  return (
    <div className='loginContainer'>
        <div className='loginBox'>
            <h1>Login to your account</h1>
            <form>
                <input type="text" placeholder='username' required/>
                <input type="text" placeholder='password' required/>
                <button type='submit'>Login</button>
            </form>
        </div>
    </div>
  );
}
