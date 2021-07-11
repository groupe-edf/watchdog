import { Reducer } from 'redux'
import {
  RepositoriesActions,
  Repository,
  RepositoryActionTypes,
  RepositoryState
} from './types'

export const initialState: RepositoryState = {
  analyzes: [],
  repositories: [],
  repository: {
    id: "",
    repository_url: "",
    visibility: ""
  }
}

const reducer: Reducer<RepositoryState> = (state: RepositoryState = initialState, action) => {
  const { payload } = action
  switch ((action as RepositoriesActions).type) {
    case RepositoryActionTypes.REPOSITORIES_APPEND: {
      const index = state.repositories.findIndex(repository => repository.id === payload.repository.id)
      let repositories = [...state.repositories]
      let repository = payload.repository as Repository
      repository.last_analysis = payload
      if (index >= 0) {
        repositories[index] = repository
      } else {
        repositories = [...repositories, repository]
      }
      return {
        ...state,
        repositories: repositories
      }
    }
    case RepositoryActionTypes.ANALYZES_FIND_ALL:
      return { ...state, analyzes: payload }
    case RepositoryActionTypes.REPOSITORIES_DELETE_BY_ID:
      return {
        ...state,
        repositories: state.repositories.filter((repository, index) => repository.id !== payload)
      }
    case RepositoryActionTypes.REPOSITORIES_FIND_ALL:
      return { ...state, repositories: payload }
    case RepositoryActionTypes.REPOSITORIES_FIND_BY_ID:
      return { ...state, repository: payload }
    default:
      return state
  }
}

export default reducer
