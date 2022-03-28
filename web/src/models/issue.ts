import { Policy } from "./policy"
import { Repository } from "./repository"

export interface Commit {
  author?: {
    date: Date,
    email: string,
    name: string,
    timezone: string
  },
  email: string,
  hash: string,
}

export interface Offender {
  object: string,
  operand: string,
  operator: string,
  value: string
}

export interface Issue {
  id: string,
  commit: Commit,
  condition_type: string,
  file?: string,
  offender?: Offender,
  policy?: Policy,
  repository?: Repository
  severity: number
}
