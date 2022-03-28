import axios from 'axios'

class GlobalService {
  async getVersion() {
    return axios.get('/version')
  }
  async evaluatePattern(data: { pattern: string, payload: string }) {
    return axios.post('/pattern', data)
  }
}

export default new GlobalService()
