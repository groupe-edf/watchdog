import { Reducer } from 'redux';
import { AuthenticationActions, UserState } from './types';

export const initialState: UserState = {
  username: '',
  password: '',
};

const reducer: Reducer<UserState> = (state: UserState = initialState, action) => {
  switch ((action as AuthenticationActions).type) {
    case '@@authentication/LOGIN':
      return { ...state, username: action.username };
    default:
      return state;
  }
}

export default reducer;
