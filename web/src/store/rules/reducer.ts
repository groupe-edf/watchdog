import { Reducer } from 'redux';
import { RuleActionTypes, RulesActions, RuleState } from './types';

export const initialState: RuleState = {
  rules: [],
};

const reducer: Reducer<RuleState> = (state: RuleState = initialState, action) => {
  const { type, payload } = action;
  switch ((action as RulesActions).type) {
    case RuleActionTypes.RULES_FIND_ALL:
      return { ...state, rules: payload };
    default:
      return state;
  }
}

export default reducer;
