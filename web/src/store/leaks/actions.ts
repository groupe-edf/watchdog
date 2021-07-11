import { ActionCreator } from 'redux';
import {
  LeakActionTypes,
  LeakFindAllAction,
  LeakFindByIdAction,
  LeakState
} from './types';

export const findAll: ActionCreator<LeakFindAllAction> = (state: LeakState) => ({
  type: LeakActionTypes.LEAKS_FIND_ALL,
  payload: {
    state
  }
})

export const findById: ActionCreator<LeakFindByIdAction> = (state: LeakState) => ({
  type: LeakActionTypes.LEAKS_FIND_BY_ID,
  payload: {
    state
  }
})
