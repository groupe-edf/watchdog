import { Action } from "redux";

export interface Integration {
  id: string
  created_at: string
  instance_name: string
  instance_url: string
  synced_at: string
  syncing_error?: string
}

export enum IntegrationActionTypes {
  INTEGRATION_FIND_ALL = "@@integrations/FIND_ALL",
  INTEGRATION_FIND_BY_ID = "@@integrations/FIND_BY_ID",
  INTEGRATION_SYNCHROONIZE = "@@integrations/SYNCHROONIZE"
}

export interface IntegrationState {
  readonly integrations: Integration[]
  readonly integration: Integration
}

export interface IntegrationsFindAllAction extends Action {
  type: IntegrationActionTypes.INTEGRATION_FIND_ALL
  payload: {
    state: IntegrationState
  }
}

export interface IntegrationsFindByIdAction extends Action {
  type: IntegrationActionTypes.INTEGRATION_FIND_BY_ID
  payload: {
    state: IntegrationState
  }
}

export interface IntegrationsSynchronizeAction extends Action {
  type: IntegrationActionTypes.INTEGRATION_SYNCHROONIZE
  payload: {
    state: IntegrationState
  }
}

export type IntegrationsActions = IntegrationsFindAllAction | IntegrationsFindByIdAction | IntegrationsSynchronizeAction
