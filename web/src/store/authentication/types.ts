import { Action } from 'redux';

export enum AuthenticationActionTypes {
  AUTHENTICATION_CURRENT_USER = "@@authentication/CURRENT_USER",
  AUTHENTICATION_LOGIN = "@@authentication/LOGIN",
  AUTHENTICATION_LOGOUT = "@@authentication/LOGOUT",
  AUTHENTICATION_REGISTER = "@@authentication/REGISTER"
}

export interface AuthenticationState {
  currentUser: {
    email: string
    first_name: string
    last_name: string
    provider: string
  }
}

export interface CurrentUserAction extends Action {
  type: AuthenticationActionTypes.AUTHENTICATION_CURRENT_USER;
  payload: {
    user: AuthenticationState;
  };
}

export interface AuthenticationLoginAction extends Action {
  type: AuthenticationActionTypes.AUTHENTICATION_LOGIN;
  payload: {
    user: AuthenticationState;
  };
}

export interface AuthenticationLogouAction extends Action {
  type: AuthenticationActionTypes.AUTHENTICATION_LOGOUT;
  payload: {
    user: AuthenticationState;
  };
}

export interface AuthenticationRegisterAction extends Action {
  type: AuthenticationActionTypes.AUTHENTICATION_REGISTER;
  payload: {
    user: AuthenticationState;
  };
}

export type AuthenticationActions = AuthenticationLoginAction | AuthenticationLogouAction | AuthenticationRegisterAction | CurrentUserAction
