import axios from 'axios'
import { Analysis, Query } from '../models'

class AnalysisService {
  async findAll(query?: Query) {
    let url = '/analyzes?'
    if (query?.conditions && query.conditions.length > 0) {
      for (let condition of query.conditions) {
        url += `&conditions=${condition.field},${condition.operator},${condition.value}`
      }
    }
    if (query?.sort && query?.sort.length > 0) {
      for (let sort of query?.sort) {
        url += `&sort=${sort.field},${sort.direction}`
      }
    }
    return axios.get<Analysis[]>(url)
  }
  async findById(id: string) {
    return axios.get<Analysis>(`/analyzes/${id}`)
  }
}

export default new AnalysisService()
