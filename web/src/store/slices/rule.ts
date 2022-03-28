import { createAsyncThunk, createSlice } from "@reduxjs/toolkit"
import { Rule } from "../../models";
import { RuleService } from "../../services";

export enum UserActionTypes {
  RULES_ADD = "@@rules/ADD",
  RULES_FIND = "@@rules/FIND",
  RULES_FIND_BY_ID = "@@rules/FIND_BY_ID"
}

export const addRule = createAsyncThunk(
  UserActionTypes.RULES_ADD,
  async (rule: Rule) => {
    const response = await RuleService.add(rule)
    return response.data
  }
)

export const getRules = createAsyncThunk(
  UserActionTypes.RULES_FIND,
  async () => {
    const response = await RuleService.findAll()
    return response.data
  }
)

export interface RuleState {
  readonly rules: Rule[]
}

export const initialState: RuleState = {
  rules: [],
}

const repositorySlice = createSlice({
  name: 'rule',
  initialState: initialState as RuleState,
  reducers: {
  },
  extraReducers: (builder) => {
    builder
      .addCase(getRules.fulfilled, (state, { payload }) => {
        state.rules = payload
      })
  }
})

const { reducer } = repositorySlice
export default reducer
