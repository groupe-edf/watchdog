import axios from 'axios'
import { User } from '../models'

class UserService {
  async findAll() {
    return axios.get<User[]>('/users')
  }
  async findById(id: number) {
    return axios.get<User>(`/users/${id}`)
  }
}

export default new UserService()
