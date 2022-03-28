import { Integration } from "./integration"
import { User } from "./user"

export interface Repository {
  id: string
  integration?: Integration
  last_analysis?: Analysis
  repository_url: string
  issues?: number
  visibility: string
}

export interface Analysis {
  id: string
  created_at: string
  created_by?: User
  duration?: number
  finished_at?: string
  last_commit_hash: string
  repository?: Repository
  started_at?: string
  state: string
  state_message?: string
  severity: string
  total_issues: number
  trigger: string
}
