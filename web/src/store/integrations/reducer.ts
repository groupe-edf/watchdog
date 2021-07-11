import { Reducer } from 'redux';
import {
  IntegrationActionTypes,
  IntegrationsActions,
  IntegrationState
} from './types';

export const initialState: IntegrationState = {
  integrations: [],
  integration: {
    id: "",
    created_at: "",
    instance_name: "",
    instance_url: "",
    synced_at: ""
  }
};

const reducer: Reducer<IntegrationState> = (state: IntegrationState = initialState, action) => {
  const { type, payload } = action;
  switch ((action as IntegrationsActions).type) {
    case IntegrationActionTypes.INTEGRATION_FIND_ALL:
      return { ...state, integrations: payload };
    case IntegrationActionTypes.INTEGRATION_FIND_BY_ID:
      return { ...state, integration: payload }
    case IntegrationActionTypes.INTEGRATION_SYNCHROONIZE:
      return { ...state, integrations: state.integrations.map(
        integration => integration.id === payload.id ? payload : integration
      )}
    default:
      return state;
  }
}

export default reducer;
