import { Reducer } from "redux";
import { JobState, JobsActions, JobActionTypes } from "./types";

export const initialState: JobState = {
  jobs: []
}

const reducer: Reducer<JobState> = (state: JobState = initialState, action) => {
  const { type, payload } = action;
  switch ((action as JobsActions).type) {
    case JobActionTypes.JOBS_FIND_ALL:
      return { ...state, jobs: payload }
    default:
      return state;
  }
}

export default reducer;
