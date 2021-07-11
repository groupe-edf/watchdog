import { ActionCreator } from 'redux'
import {
  AppendAction,
  AnalyzesFindAllAction,
  RepositoriesDeleteByIdAction,
  RepositoriesFindAllAction,
  RepositoriesFindByIdAction,
  RepositoryActionTypes,
  RepositoryState
} from './types'

export const append: ActionCreator<AppendAction> = (state: RepositoryState) => ({
  type: RepositoryActionTypes.REPOSITORIES_APPEND,
  payload: {
    state,
  }
})

export const findAllAnalyzes: ActionCreator<AnalyzesFindAllAction> = (state: RepositoryState) => ({
  type: RepositoryActionTypes.ANALYZES_FIND_ALL,
  payload: {
    state,
  }
})

export const deleteRepository: ActionCreator<RepositoriesDeleteByIdAction> = (state: RepositoryState) => ({
  type: RepositoryActionTypes.REPOSITORIES_DELETE_BY_ID,
  payload: {
    state,
  }
})

export const findAll: ActionCreator<RepositoriesFindAllAction> = (state: RepositoryState) => ({
  type: RepositoryActionTypes.REPOSITORIES_FIND_ALL,
  payload: {
    state,
  }
})

export const findById: ActionCreator<RepositoriesFindByIdAction> = (state: RepositoryState) => ({
  type: RepositoryActionTypes.REPOSITORIES_FIND_BY_ID,
  payload: {
    state
  }
})
