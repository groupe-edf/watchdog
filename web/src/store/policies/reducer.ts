import { Reducer } from 'redux';
import { PolicyActionTypes, PoliciesActions, PolicyState } from './types';

export const initialState: PolicyState = {
  policies: [],
  policy: {
    id: 0,
    display_name: "",
    enabled: false,
    type: ""
  }
};

const reducer: Reducer<PolicyState> = (state: PolicyState = initialState, action) => {
  const { type, payload } = action;
  switch ((action as PoliciesActions).type) {
    case PolicyActionTypes.POLICIES_FIND_ALL:
      return { ...state, policies: payload }
    case PolicyActionTypes.POLICIES_FIND_BY_ID:
      return { ...state, policy: payload }
    default:
      return state;
  }
}

export default reducer;
