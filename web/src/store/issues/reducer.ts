import { Reducer } from 'redux';
import { IssueActionTypes, IssuesActions, IssueState } from './types';

export const initialState: IssueState = {
  issues: [],
};

const reducer: Reducer<IssueState> = (state: IssueState = initialState, action) => {
  const { type, payload } = action;
  switch ((action as IssuesActions).type) {
    case IssueActionTypes.ISSUES_FIND_ALL:
      return { ...state, issues: payload };
    default:
      return state;
  }
}

export default reducer;
