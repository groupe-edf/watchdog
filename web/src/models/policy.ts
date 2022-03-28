export interface Policy {
  id: number
  conditions?: PolicyCondition[]
  description?: string
  display_name: string
  enabled: boolean
  severity: string
  type: string
}

export interface PolicyCondition {
  id: number
  pattern: string
  policy_id: number
  type: string
}
