import { Reducer } from 'redux';
import { AuthenticationActions, AuthenticationActionTypes, AuthenticationState } from './types';

export const initialState: AuthenticationState = {
  currentUser: {
    email: '',
    first_name: '',
    last_name: '',
    provider: 'local'
  }
};

const reducer: Reducer<AuthenticationState> = (state: AuthenticationState = initialState, action) => {
  const { type, payload } = action;
  switch ((action as AuthenticationActions).type) {
    case AuthenticationActionTypes.AUTHENTICATION_CURRENT_USER:
      return { ...state, currentUser: payload };
    case AuthenticationActionTypes.AUTHENTICATION_LOGIN:
      return { ...state, currentUser: payload };
    case AuthenticationActionTypes.AUTHENTICATION_LOGOUT:
      return { ...state, currentUser: {} };
    case AuthenticationActionTypes.AUTHENTICATION_REGISTER:
      return { ...state, currentUser: payload };
    default:
      return state;
  }
}

export default reducer;
