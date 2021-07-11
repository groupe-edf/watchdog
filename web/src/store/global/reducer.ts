import { Reducer } from "redux";
import { GlobalActions, GlobalActionTypes, GlobalState } from "./types";

export const initialState: GlobalState = {
  api_keys: {},
  current_api_key: {},
  version: {},
  settings: {}
};

const reducer: Reducer<GlobalState> = (state: GlobalState = initialState, action) => {
  const { type, payload } = action;
  switch ((action as GlobalActions).type) {
    case GlobalActionTypes.GLOBAL_API_KEYS:
      return { ...state, api_keys: payload }
    case GlobalActionTypes.GLOBAL_CURRENT_API_KEY:
      return { ...state, current_api_key: payload }
    case GlobalActionTypes.GLOBAL_VERSION:
      return { ...state, version: payload }
    case GlobalActionTypes.GLOBAL_SETTINGS:
      return { ...state, settings: payload }
    default:
      return state
  }
}

export default reducer
