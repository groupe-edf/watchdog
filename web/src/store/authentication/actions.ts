import { ActionCreator } from 'redux';
import {
  AuthenticationActionTypes,
  AuthenticationLoginAction,
  AuthenticationLogouAction,
  AuthenticationRegisterAction,
  AuthenticationState,
  CurrentUserAction
} from './types';

export const currentUser: ActionCreator<CurrentUserAction> = (user: AuthenticationState) => ({
  type: AuthenticationActionTypes.AUTHENTICATION_CURRENT_USER,
  payload: {
    user,
  },
});

export const login: ActionCreator<AuthenticationLoginAction> = (user: AuthenticationState) => ({
  type: AuthenticationActionTypes.AUTHENTICATION_LOGIN,
  payload: {
    user,
  },
});

export const logout: ActionCreator<AuthenticationLogouAction> = (user: AuthenticationState) => ({
  type: AuthenticationActionTypes.AUTHENTICATION_LOGOUT,
  payload: {
    user,
  },
});

export const register: ActionCreator<AuthenticationRegisterAction> = (user: AuthenticationState) => ({
  type: AuthenticationActionTypes.AUTHENTICATION_REGISTER,
  payload: {
    user,
  },
});
