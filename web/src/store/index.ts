import { combineReducers } from '@reduxjs/toolkit'
import analysisReducer from './slices/analysis'
import authenticationReducer from './slices/authentication'
import globalReducer from './slices/global'
import integrationReducer from './slices/integration'
import issueReducer from './slices/issue'
import leakReducer from './slices/leak'
import notificationReducer from './slices/notification'
import policyReducer from './slices/policy'
import repositoryReducer from './slices/repository'
import ruleReducer from './slices/rule'
import userReducer from './slices/user'

const rootReducer = combineReducers({
  analyzes: analysisReducer,
  authentication: authenticationReducer,
  global: globalReducer,
  integrations: integrationReducer,
  issues: issueReducer,
  leaks: leakReducer,
  notifications: notificationReducer,
  policies: policyReducer,
  repositories: repositoryReducer,
  rules: ruleReducer,
  users: userReducer
})
export default rootReducer
