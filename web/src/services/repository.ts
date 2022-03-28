import axios from 'axios'
import { Analysis, Query, Repository } from '../models'

class RepositoryService {
  async analyze(data: any) {
    return axios.post<Repository[]>('/analyze', data)
  }
  async analyzeById(id: string, data: any) {
    return axios.post<Analysis>(`/repositories/${id}/analyze`, data)
  }
  async deleteById(id: string) {
    return axios.delete<Repository[]>(`/repositories/${id}`)
  }
  async findAll(query?: any) {
    return axios.get<Repository[]>(`/repositories?${new URLSearchParams(query).toString()}`)
  }
  async findById(id: string) {
    return axios.get<Repository>(`/repositories/${id}`)
  }
  async getBadge(id: string) {
    return axios.get(`/repositories/${id}/badge`)
  }
}

export default new RepositoryService()
