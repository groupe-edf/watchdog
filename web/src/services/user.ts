import { API_PATH } from "../constants";
import { User } from "../store/users/types";
import { Query, fetchData } from "./commons";

class UserService {
  async changePassword(data: {
    current_password: string,
    password: string,
    confirm_password: string
  }) {
    return fetchData("PUT", `${API_PATH}/password`, data)
  }
  async findAll() {
    return fetchData<User[]>("GET", `${API_PATH}/users`)
  }
  async findById(id: number) {
    return fetchData<User>("GET", `${API_PATH}/users/${id}`)
  }
}

export default new UserService()
