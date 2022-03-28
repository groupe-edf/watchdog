import { createAsyncThunk, createSlice } from "@reduxjs/toolkit"
import { Integration } from "../../models/integration"
import { IntegrationService } from '../../services'

export enum IntegrationActionTypes {
  INTEGRATIONS_ADD = "@@integrations/ADD",
  INTEGRATIONS_DELETE = "@@integrations/DELETE",
  INTEGRATIONS_FIND = "@@integrations/FIND",
  INTEGRATIONS_FIND_BY_ID = "@@integrations/FIND_BY_ID",
  INTEGRATIONS_SYNCHRONIZE = "@@integrations/SYNCHRONIZE"
}

export const addIntegration = createAsyncThunk(
  IntegrationActionTypes.INTEGRATIONS_ADD,
  async (data: any, thunkAPI) => {
    try {
      const response = await IntegrationService.create(data)
      thunkAPI.dispatch(getIntegrations())
      return response.data
    } catch (error) {
      return thunkAPI.rejectWithValue(error)
    }
  }
)

export const deleteIntegration = createAsyncThunk(
  IntegrationActionTypes.INTEGRATIONS_DELETE,
  async (integrationId: string, thunkAPI) => {
    const response = await IntegrationService.delete(integrationId)
    thunkAPI.dispatch(getIntegrations())
    return response.data
  }
)

export const getIntegration = createAsyncThunk(
  IntegrationActionTypes.INTEGRATIONS_FIND_BY_ID,
  async (integrationId: string, thunkAPI) => {
    const response = await IntegrationService.findById(integrationId)
    return response.data
  }
)

export const getIntegrations = createAsyncThunk(
  IntegrationActionTypes.INTEGRATIONS_FIND,
  async () => {
    const response = await IntegrationService.findAll()
    return response.data
  }
)

export const synchronizeInstance = createAsyncThunk(
  IntegrationActionTypes.INTEGRATIONS_SYNCHRONIZE,
  async (integrationId: string, thunkAPI) => {
    try {
      const response = await IntegrationService.synchronize(integrationId)
      thunkAPI.dispatch(getIntegrations())
      return response.data
    } catch (error) {
      return thunkAPI.rejectWithValue(error)
    }
  }
)

export interface IntegrationState {
  readonly integration: Integration
  readonly integrations: Integration[]
}

export const initialState: IntegrationState = {
  integration: {
    id: '',
    created_at: '',
    instance_name: '',
    instance_url: ''
  },
  integrations: []
}

const integrationSlice = createSlice({
  name: 'integration',
  initialState: initialState as IntegrationState,
  reducers: {
  },
  extraReducers: (builder) => {
    builder
      .addCase(getIntegration.fulfilled, (state, { payload }) => {
        state.integration = payload
      })
      .addCase(getIntegrations.fulfilled, (state, { payload }) => {
        state.integrations = payload
      })
  }
})

const { reducer } = integrationSlice
export default reducer
