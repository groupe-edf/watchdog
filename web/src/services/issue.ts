import axios from 'axios'
import { Issue } from '../models'

class IssueService {
  async findAll(query?: any) {
    return axios.get<Issue[]>(`/issues?${new URLSearchParams(query).toString()}`)
  }
}

export default new IssueService()
