import { ActionCreator } from 'redux';
import {
  JobFindAllAction,
  JobState,
  JobActionTypes
} from './types';

export const findAll: ActionCreator<JobFindAllAction> = (state: JobState) => ({
  type: JobActionTypes.JOBS_FIND_ALL,
  payload: {
    state
  }
})
