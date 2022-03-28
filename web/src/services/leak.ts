import axios from 'axios'
import { Leak } from '../models'

class LeakService {
  async findAll(query?: any) {
    return axios.get<Leak[]>(`/leaks?${new URLSearchParams(query).toString()}`)
  }
  async findById(id: string) {
    return axios.get<Leak>(`/leaks/${id}`)
  }
}

export default new LeakService()
