export interface Rule {
  id: number
  display_name: string
  enabled: boolean
  pattern?: string
  severity: string
  tags: string[]
}
