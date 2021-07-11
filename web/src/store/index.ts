import { combineReducers } from 'redux'
import { connectRouter, RouterState } from 'connected-react-router'
import { History } from 'history'

import authenticationReducer from './authentication/reducer'
import { AuthenticationState } from './authentication/types'
import categoriesReducer from './categories/reducer'
import { CategoryState } from './categories/types'
import globalReducer from './global/reducer'
import { GlobalState } from './global/types'
import integrationsReducer from './integrations/reducer'
import { IntegrationState } from './integrations/types'
import issuesReducer from './issues/reducer'
import { IssueState } from './issues/types'
import leaksReducer from './leaks/reducer'
import { LeakState } from './leaks/types'
import policiesReducer from './policies/reducer'
import { PolicyState } from './policies/types'
import repositoriesReducer from './repositories/reducer'
import { RepositoryState } from './repositories/types'
import rulesReducer from './rules/reducer'
import { RuleState } from './rules/types'
import usersReducer from './users/reducer'
import { UserState } from './users/types'
import { JobState } from './jobs/types'
import jobsReducer from './jobs/reducer'

export interface ApplicationState {
  authentication: AuthenticationState,
  categories: CategoryState,
  global: GlobalState,
  integrations: IntegrationState,
  issues: IssueState,
  jobs: JobState,
  leaks: LeakState,
  policies: PolicyState,
  repositories: RepositoryState,
  rules: RuleState,
  router: RouterState,
  users: UserState
}

export const createRootReducer = (history: History) => combineReducers({
  authentication: authenticationReducer,
  categories: categoriesReducer,
  global: globalReducer,
  integrations: integrationsReducer,
  issues: issuesReducer,
  jobs: jobsReducer,
  leaks: leaksReducer,
  policies: policiesReducer,
  repositories: repositoriesReducer,
  rules: rulesReducer,
  router: connectRouter(history),
  users: usersReducer
})
