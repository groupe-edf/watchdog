import { Action } from "redux";

export interface Policy {
  id: number,
  conditions?: PolicyCondition[],
  description?: string,
  display_name: string,
  enabled: boolean,
  type: string
}

export interface PolicyCondition {
  type: string,
  pattern: string
}

export enum PolicyActionTypes {
  POLICIES_FIND_ALL = "@@policies/FIND_ALL",
  POLICIES_FIND_BY_ID = "@@policies/FIND_BY_ID"
}

export interface PolicyState {
  policies: Policy[];
  policy: Policy
}

export interface PolicyFindAllAction extends Action {
  type: PolicyActionTypes.POLICIES_FIND_ALL;
  payload: {
    state: PolicyState;
  };
}

export interface PolicyFindByIdAction extends Action {
  type: PolicyActionTypes.POLICIES_FIND_BY_ID;
  payload: {
    state: PolicyState;
  };
}

export type PoliciesActions = PolicyFindAllAction | PolicyFindByIdAction
