import { ActionCreator } from 'redux';
import {
  PolicyActionTypes,
  PolicyFindAllAction,
  PolicyFindByIdAction,
  PolicyState
} from './types';

export const findAll: ActionCreator<PolicyFindAllAction> = (state: PolicyState) => ({
  type: PolicyActionTypes.POLICIES_FIND_ALL,
  payload: {
    state
  }
})

export const findById: ActionCreator<PolicyFindByIdAction> = (state: PolicyState) => ({
  type: PolicyActionTypes.POLICIES_FIND_BY_ID,
  payload: {
    state
  }
})
