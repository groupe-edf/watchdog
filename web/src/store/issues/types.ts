import { Action } from "redux";
import { Policy } from "../policies/types";
import { Repository } from "../repositories/types";

export interface Commit {
  author: string,
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
  policy: Policy,
  repository: Repository
  severity: number,
}

export enum IssueActionTypes {
  ISSUES_FIND_ALL = "@@issues/FIND_ALL"
}

export interface IssueState {
  issues: Issue[];
}

export interface IssueFindAllAction extends Action {
  type: IssueActionTypes.ISSUES_FIND_ALL;
  payload: {
    state: IssueState;
  };
}

export type IssuesActions = IssueFindAllAction;
