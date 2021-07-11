import { ActionCreator } from 'redux';
import {
  UserActionTypes,
  UserFindAllAction,
  UserState
} from './types';

export const findAll: ActionCreator<UserFindAllAction> = (state: UserState) => ({
  type: UserActionTypes.USERS_FIND_ALL,
  payload: {
    state,
  },
});
