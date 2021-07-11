import { Reducer, combineReducers } from 'redux'
import authenticationReducer from './authentication/reducer'
import { UserState } from './authentication/types';

export interface ApplicationState {
  authentication: UserState
}

export const reducers: Reducer<ApplicationState> = combineReducers<ApplicationState>({
  authentication: authenticationReducer
});
