import { Repository } from "./repository"
import { Rule } from "./rule"

export interface Leak {
  id: string,
  author_email: string,
  author_name: string,
  commit_hash: string,
  created_at: string,
  file: string,
  line: string,
  line_number: number,
  occurence: number,
  offender: string,
  repository: Repository,
  rule: Rule,
  severity: string
}
