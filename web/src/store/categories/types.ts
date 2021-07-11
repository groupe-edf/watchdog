import { Action } from "redux";

export interface Category {
  id: string
  description?: string
  extension: string
  level: string
  left: string
  right: string
  title: string
  value: string
}

export enum CategoryActionTypes {
  CATEGORY_FIND_ALL = "@@categories/FIND_ALL"
}

export interface CategoryState {
  readonly categories: Category[]
}

export interface CategoriesFindAllAction extends Action {
  type: CategoryActionTypes.CATEGORY_FIND_ALL
  payload: {
    state: CategoryState
  }
}

export type CategoriesActions = CategoriesFindAllAction
