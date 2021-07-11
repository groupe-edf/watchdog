import { ActionCreator } from 'redux'
import {
  GlobalSettingsAction,
  GlobalVersionAction,
  GlobalActionTypes,
  GlobalState,
  GlobalAPIKeysAction,
  GlobalCurrentAPIKeyAction
} from './types'

export const getAPIKeys: ActionCreator<GlobalAPIKeysAction> = (state: GlobalState) => ({
  type: GlobalActionTypes.GLOBAL_API_KEYS,
  payload: {
    state
  }
})

export const getCurrentAPIKey: ActionCreator<GlobalCurrentAPIKeyAction> = (state: GlobalState) => ({
  type: GlobalActionTypes.GLOBAL_CURRENT_API_KEY,
  payload: {
    state
  }
})

export const getSettings: ActionCreator<GlobalSettingsAction> = (state: GlobalState) => ({
  type: GlobalActionTypes.GLOBAL_SETTINGS,
  payload: {
    state
  }
})

export const getVersion: ActionCreator<GlobalVersionAction> = (state: GlobalState) => ({
  type: GlobalActionTypes.GLOBAL_VERSION,
  payload: {
    state
  }
})
