import { Action } from "redux";

export interface User {
  id: number
  created_at: string
  email: string
  first_name: string
  last_login: string
  last_name: string
  locked: boolean
  provider: string
  role: string
  username: string
}

export enum UserActionTypes {
  USERS_FIND_ALL = "@@users/FIND_ALL"
}

export interface UserState {
  users: User[];
}

export interface UserFindAllAction extends Action {
  type: UserActionTypes.USERS_FIND_ALL;
  payload: {
    state: UserState;
  };
}

export type UsersActions = UserFindAllAction;
