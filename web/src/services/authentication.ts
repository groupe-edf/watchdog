import { API_PATH } from '../constants';
import { Query, fetchData } from "./commons";

class AuthenticationService {
  async login({ email, password } : { email: string, password: string }) {
    return fetchData("POST", `${API_PATH}/login`, {
      email: email,
      password: password
    });
  }
  logout() {
    localStorage.removeItem('user');
  }
  async register({ email, first_name, last_name, password } : { email: string, first_name: string, last_name: string, password: string }) {
    return fetchData("POST", `${API_PATH}/register`, {
      email: email,
      first_name: first_name,
      last_name: last_name,
      password: password
    });
  }
}

export default new AuthenticationService();
