import { ActionCreator } from 'redux';
import {
  CategoriesFindAllAction,
  CategoryActionTypes,
  CategoryState
} from './types';

export const findAll: ActionCreator<CategoriesFindAllAction> = (state: CategoryState) => ({
  type: CategoryActionTypes.CATEGORY_FIND_ALL,
  payload: {
    state
  }
})
