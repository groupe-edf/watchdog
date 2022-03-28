import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { Issue, Query } from "../../models";
import { IssueService } from "../../services";

export enum IssueActionTypes {
  ISSUES_FIND = "@@issues/FIND",
  ISSUES_FIND_BY_ID = "@@issues/FIND_BY_ID",
}

export const getIssues = createAsyncThunk(
  IssueActionTypes.ISSUES_FIND,
  async (query: Query, thunkAPI) => {
    const response = await IssueService.findAll(query)
    return response.data
  }
)

export interface IssueState {
  readonly issue: Issue
  readonly issues: Issue[]
}

export const initialState: IssueState = {
  issue: {
    id: "",
    commit: {
      email: "",
      hash: ""
    },
    condition_type: "",
    severity: 0
  },
  issues: [],
}

const issueSlice = createSlice({
  name: 'issue',
  initialState: initialState as IssueState,
  reducers: {
  },
  extraReducers: (builder) => {
    builder
      .addCase(getIssues.fulfilled, (state, { payload }) => {
        state.issues = payload
      })
  }
})

const { reducer } = issueSlice
export default reducer
