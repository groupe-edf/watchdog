import { ActionCreator } from 'redux';
import {
  AuthenticationLoginAction,
  UserState
} from './types';

export const login: ActionCreator<AuthenticationLoginAction> = (user: UserState) => ({
  type: '@@authentication/LOGIN',
  payload: {
    user,
  },
});
