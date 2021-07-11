import { ActionCreator } from 'redux';
import {
  RuleActionTypes,
  RuleFindAllAction,
  RuleState
} from './types';

export const findAll: ActionCreator<RuleFindAllAction> = (state: RuleState) => ({
  type: RuleActionTypes.RULES_FIND_ALL,
  payload: {
    state,
  },
});
