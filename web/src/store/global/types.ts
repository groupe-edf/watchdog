import { Action } from 'redux';
import { Query } from '../../services/commons';

export enum GlobalActionTypes {
  GLOBAL_API_KEYS = "@@global/API_KEYS",
  GLOBAL_CURRENT_API_KEY = "@@global/CURRENT_API_KEY",
  GLOBAL_SETTINGS = "@@global/SETTINGS",
  GLOBAL_VERSION = "@@global/VERSION",
}

export interface AccessToken {

}

export interface GlobalState {
  api_keys: any
  current_api_key: any
  settings: any
  version: any
}

export interface TableState {
  isLoading: boolean
  query: Query
  totalItems: number
}

export interface GlobalAPIKeysAction extends Action {
  type: GlobalActionTypes.GLOBAL_API_KEYS
  payload: {
    state: GlobalState
  }
}

export interface GlobalCurrentAPIKeyAction extends Action {
  type: GlobalActionTypes.GLOBAL_CURRENT_API_KEY
  payload: {
    state: GlobalState
  }
}

export interface GlobalSettingsAction extends Action {
  type: GlobalActionTypes.GLOBAL_SETTINGS
  payload: {
    state: GlobalState
  }
}

export interface GlobalVersionAction extends Action {
  type: GlobalActionTypes.GLOBAL_VERSION
  payload: {
    state: GlobalState
  }
}

export type GlobalActions = GlobalAPIKeysAction | GlobalCurrentAPIKeyAction | GlobalSettingsAction | GlobalVersionAction
