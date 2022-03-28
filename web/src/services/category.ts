import axios from 'axios'
import { Category, Query } from '../models'

class CategoryService {
  async findAll(query?: Query) {
    let url = '/categories'
    if (query?.conditions && query.conditions.length > 0) {
      for (let condition of query.conditions) {
        url += `&conditions=${condition.field},${condition.operator},${condition.value}`
      }
    }
    return axios.get<Category[]>(url)
  }
}

export default new CategoryService()
