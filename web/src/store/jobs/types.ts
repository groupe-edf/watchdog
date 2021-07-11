import { Action } from "redux"

export interface Job {
  id: string
  error_count: number
  priority: string
  queue: string
  started_at: string
  type: string
}

export enum JobActionTypes {
  JOBS_FIND_ALL = "@@jobs/FIND_ALL"
}

export interface JobState {
  jobs: Job[]
}

export interface JobFindAllAction extends Action {
  type: JobActionTypes.JOBS_FIND_ALL;
  payload: {
    state: JobState
  }
}

export type JobsActions = JobFindAllAction
