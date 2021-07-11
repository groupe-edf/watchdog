import { Action } from "redux";

export interface Rule {
  id: number,
  display_name: string,
  enabled: boolean,
  severity: string,
  tags: string[],
}

export enum RuleActionTypes {
  RULES_FIND_ALL = "@@rules/FIND_ALL"
}

export interface RuleState {
  rules: Rule[];
}

export interface RuleFindAllAction extends Action {
  type: RuleActionTypes.RULES_FIND_ALL;
  payload: {
    state: RuleState;
  };
}

export type RulesActions = RuleFindAllAction;
