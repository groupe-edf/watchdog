export interface Condition {
  field: string
  operator: string
  value?: string
}

export interface Pagination {
  currentPage?: number
  pagesToShow?: number
  itemsPerPage: number
  offset: number
  totalItems: number
}

export class Query {
  limit?: number
  query?: string
  conditions?: Condition[]
  offset?: number
  sort?: Sort[]
  toString(): string {
    let query = ""
    return query
  }
}

export interface Result {
  data: any
  error: any
  pagination: Range
}

export interface Sort {
  field: string
  direction: string
}

export interface Version {
  platform: string
  version: string
}
