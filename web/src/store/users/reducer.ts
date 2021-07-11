import { Reducer } from 'redux';
import { UserActionTypes, UsersActions, UserState } from './types';

export const initialState: UserState = {
  users: [],
};

const reducer: Reducer<UserState> = (state: UserState = initialState, action) => {
  const { type, payload } = action;
  switch ((action as UsersActions).type) {
    case UserActionTypes.USERS_FIND_ALL:
      return { ...state, users: payload };
    default:
      return state;
  }
}

export default reducer;
