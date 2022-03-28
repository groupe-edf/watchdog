import axios from 'axios'

class IntegrationService {
  async create(data: any) {
    return axios.post('/integrations', data)
  }
  async delete(id: string) {
    return axios.delete(`/integrations/${id}`)
  }
  async findAll() {
    return axios.get('/integrations')
  }
  async findById(id: string) {
    return axios.get(`/integrations/${id}`)
  }
  async getGroups(id: string) {
    return axios.get(`/integrations/${id}/groups`)
  }
  async installWebhook(data: any) {
    return axios.post(`/integrations/${data.intergration_id}/webhooks`, data)
  }
  async synchronize(id: string) {
    return axios.get(`/integrations/${id}/synchronize`)
  }
}

export default new IntegrationService()
