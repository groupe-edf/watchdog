import axios from 'axios'

class AuthenticationService {
  getProfile() {
    return axios.get('/profile')
  }
  login({ email, password } : { email: string, password: string }) {
    return axios.post('/authentication/login', {
      email,
      password
    })
  }
  logout() {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }
  refreshToken() {
  }
  register({ email, first_name, last_name, password } : { email: string, first_name: string, last_name: string, password: string }) {
    return axios.post('/authentication/register', {
      email: email,
      first_name: first_name,
      last_name: last_name,
      password: password
    });
  }
}

export default new AuthenticationService()
