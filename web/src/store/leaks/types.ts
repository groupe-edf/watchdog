import { Action } from "redux"
import { Repository } from "../repositories/types"
import { Rule } from "../rules/types"

export interface Leak {
  id: string,
  author: string,
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

export enum LeakActionTypes {
  LEAKS_FIND_ALL = "@@leaks/FIND_ALL",
  LEAKS_FIND_BY_ID = "@@leaks/FIND_BY_ID"
}

export interface LeakState {
  leaks: Leak[]
  leak: Leak
}

export interface LeakFindAllAction extends Action {
  type: LeakActionTypes.LEAKS_FIND_ALL;
  payload: {
    state: LeakState
  }
}

export interface LeakFindByIdAction extends Action {
  type: LeakActionTypes.LEAKS_FIND_BY_ID;
  payload: {
    state: LeakState;
  };
}

export type LeaksActions = LeakFindAllAction | LeakFindByIdAction
