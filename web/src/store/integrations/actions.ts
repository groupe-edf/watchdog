import { ActionCreator } from 'redux';
import {
  IntegrationsFindAllAction,
  IntegrationsFindByIdAction,
  IntegrationsSynchronizeAction,
  IntegrationActionTypes,
  IntegrationState
} from './types';

export const findAll: ActionCreator<IntegrationsFindAllAction> = (state: IntegrationState) => ({
  type: IntegrationActionTypes.INTEGRATION_FIND_ALL,
  payload: {
    state
  }
})

export const findById: ActionCreator<IntegrationsFindByIdAction> = (state: IntegrationState) => ({
  type: IntegrationActionTypes.INTEGRATION_FIND_BY_ID,
  payload: {
    state
  }
})

export const synchronize: ActionCreator<IntegrationsSynchronizeAction> = (state: IntegrationState) => ({
  type: IntegrationActionTypes.INTEGRATION_SYNCHROONIZE,
  payload: {
    state
  }
})
