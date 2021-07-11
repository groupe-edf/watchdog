import { ActionCreator } from 'redux';
import {
  IssueActionTypes,
  IssueFindAllAction,
  IssueState
} from './types';

export const findAll: ActionCreator<IssueFindAllAction> = (state: IssueState) => ({
  type: IssueActionTypes.ISSUES_FIND_ALL,
  payload: {
    state,
  },
});
