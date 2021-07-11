import { Reducer } from 'redux';
import {
  CategoryActionTypes,
  CategoriesActions,
  CategoryState
} from './types';

export const initialState: CategoryState = {
  categories: []
}

const reducer: Reducer<CategoryState> = (state: CategoryState = initialState, action) => {
  const { type, payload } = action;
  switch ((action as CategoriesActions).type) {
    case CategoryActionTypes.CATEGORY_FIND_ALL:
      return { ...state, categories: payload };
    default:
      return state;
  }
}

export default reducer;
