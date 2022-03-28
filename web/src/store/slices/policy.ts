import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { Policy, PolicyCondition } from "../../models"
import { PolicyService } from "../../services"

export enum PolicyActionTypes {
  POLICIES_ADD = "@@policies/ADD",
  POLICIES_ADD_CONDITION = "@@policies/ADD_CONDITION",
  POLICIES_DELETE_CONDITION = "@@policies/DELETE_CONDITION",
  POLICIES_DELETE_POLICY = "@@policies/DELETE_POLICY",
  POLICIES_FIND = "@@policies/FIND",
  POLICIES_FIND_BY_ID = "@@policies/FIND_BY_ID",
  POLICIES_TOGGLE = "@@policies/TOGGLE"
}

export const addPolicy = createAsyncThunk(
  PolicyActionTypes.POLICIES_ADD,
  async (policy: Policy, thunkAPI) => {
    const response = await PolicyService.create(policy)
    thunkAPI.dispatch(getPolicies())
    return response.data
  }
)

export const addPolicyCondition = createAsyncThunk(
  PolicyActionTypes.POLICIES_ADD_CONDITION,
  async (condition: PolicyCondition, thunkAPI) => {
    const response = await PolicyService.addCondition(condition)
    return response.data
  }
)

export const deleteCondition = createAsyncThunk(
  PolicyActionTypes.POLICIES_DELETE_CONDITION,
  async (data: { policyId: string, condition: PolicyCondition }, thunkAPI) => {
    const response = await PolicyService.deleteCondtion(data.policyId, data.condition)
    thunkAPI.dispatch(getPolicy(Number(data.policyId)))
    return response.data
  }
)

export const deletePolicy = createAsyncThunk(
  PolicyActionTypes.POLICIES_DELETE_POLICY,
  async (policy: Policy, thunkAPI) => {
    const response = await PolicyService.delete(policy.id)
    return response.data
  }
)

export const getPolicies = createAsyncThunk(
  PolicyActionTypes.POLICIES_FIND,
  async () => {
    const response = await PolicyService.findAll()
    return response.data
  }
)

export const getPolicy = createAsyncThunk(
  PolicyActionTypes.POLICIES_FIND_BY_ID,
  async (policyId: number) => {
    const response = await PolicyService.findById(policyId)
    return response.data
  }
)

export const togglePolicy = createAsyncThunk(
  PolicyActionTypes.POLICIES_TOGGLE,
  async (policy: Policy) => {
    await PolicyService.toggle(policy)
    return {
      id: policy.id,
      enabled: !policy.enabled
    }
  }
)

export interface PolicyState {
  readonly policy: Policy
  readonly policies: Policy[]
}

export const initialState: PolicyState = {
  policy: {
    id: 0,
    display_name: "",
    enabled: false,
    severity: "",
    type: ""
  },
  policies: [],
}

const policySlice = createSlice({
  name: 'policy',
  initialState: initialState as PolicyState,
  reducers: {
  },
  extraReducers: (builder) => {
    builder
      .addCase(addPolicyCondition.fulfilled, (state, { payload }) => {
        state.policy.conditions?.push(payload)
      })
      .addCase(deletePolicy.fulfilled, (state, { payload }) => {
        state.policy = initialState.policy
      })
      .addCase(getPolicy.fulfilled, (state, { payload }) => {
        state.policy = payload
      })
      .addCase(getPolicies.fulfilled, (state, { payload }) => {
        state.policies = payload
      })
      .addCase(togglePolicy.fulfilled, (state, { payload }) => {
        if (state.policy && state.policy.id === payload.id) {
          state.policy.enabled = payload.enabled
        }
        const index = state.policies.findIndex(policy => policy.id === payload.id)
        if (index !== -1) {
          state.policies[index].enabled = payload.enabled
        }
      })
  }
})

const { reducer } = policySlice
export default reducer
