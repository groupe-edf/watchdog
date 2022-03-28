import axios from 'axios'
import { Rule } from '../models'

class RuleService {
  async add(rule: Rule) {
    return axios.post<Rule>('/rules', rule)
  }
  async findAll() {
    return axios.get<Rule[]>('/rules')
  }
  async findById(id: number) {
    return axios.get<Rule>(`/rules/${id}`)
  }
}

export default new RuleService()
