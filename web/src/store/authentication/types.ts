import { Action } from 'redux';

export interface UserState {
  username: string;
  password: string;
}

export interface AuthenticationLoginAction extends Action {
  type: '@@authentication/LOGIN';
  payload: {
    user: UserState;
  };
}

export type AuthenticationActions = AuthenticationLoginAction;
