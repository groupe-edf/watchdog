import { Reducer } from 'redux';
import { LeakActionTypes, LeaksActions, LeakState } from './types';

export const initialState: LeakState = {
  leaks: [],
  leak: {
    id: "",
    author: "",
    commit_hash: "",
    created_at: "",
    file: "",
    line: "",
    line_number: 0,
    occurence: 0,
    offender: "",
    repository: {
      id: "",
      repository_url: "",
      visibility: ""
    },
    rule: {
      id: 0,
      display_name: "",
      enabled: true,
      severity: "",
      tags: [],
    },
    severity: ""
  }
};

const reducer: Reducer<LeakState> = (state: LeakState = initialState, action) => {
  const { type, payload } = action;
  switch ((action as LeaksActions).type) {
    case LeakActionTypes.LEAKS_FIND_ALL:
      return { ...state, leaks: payload }
    case LeakActionTypes.LEAKS_FIND_BY_ID:
      return { ...state, leak: payload }
    default:
      return state;
  }
}

export default reducer;
