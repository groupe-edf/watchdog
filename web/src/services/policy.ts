import axios from 'axios'
import { Policy, PolicyCondition } from '../models'

class PolicyService {
  async addCondition(condition: PolicyCondition) {
    return axios.post<PolicyCondition>(`/policies/${condition.policy_id}/conditions`, condition)
  }
  async create(policy: Policy) {
    return axios.post('/policies', policy)
  }
  async delete(policyId: number) {
    return axios.delete(`/policies/${policyId}`)
  }
  async deleteCondtion(policyId: string, condition: PolicyCondition) {
    return axios.delete(`/policies/${policyId}/conditions/${condition.id}`)
  }
  async findAll() {
    return axios.get<Policy[]>('/policies')
  }
  async findById(id: number) {
    return axios.get<Policy>(`/policies/${id}`)
  }
  async toggle(policy: Policy) {
    return axios.put<Policy>(`/policies/${policy.id}/toggle`, {
      enabled: !policy.enabled
    })
  }
  async update(policy: Policy) {
    return axios.put<Policy>(`/policies/${policy.id}`, policy)
  }
}

export default new PolicyService()
