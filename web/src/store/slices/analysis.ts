import { createAsyncThunk, createSlice } from "@reduxjs/toolkit"
import { Analysis, Query } from "../../models"
import { AnalysisService } from "../../services"

export enum AnalysisActionTypes {
  ANALYZES_FIND = "@@analyzes/FIND",
  ANALYZES_FIND_BY_ID = "@@analyzes/FIND_BY_ID",
}

export const getAnalysisById = createAsyncThunk(
  AnalysisActionTypes.ANALYZES_FIND_BY_ID,
  async (repositoryId: string, thunkAPI) => {
    try {
      const response = await AnalysisService.findById(repositoryId)
      return response.data
    } catch (error) {
      return thunkAPI.rejectWithValue(error)
    }
  }
)

export const getAnalyzes = createAsyncThunk(
  AnalysisActionTypes.ANALYZES_FIND,
  async (query: Query, thunkAPI) => {
    try {
      const response = await AnalysisService.findAll(query)
      return response.data
    } catch (error) {
      return thunkAPI.rejectWithValue(error)
    }
  }
)

export interface AnalysisState {
  readonly analysis: Analysis
  readonly analyzes: Analysis[]
}

export const initialState: AnalysisState = {
  analysis: {
    id: "",
    created_at: "",
    last_commit_hash: "",
    state: "",
    severity: "",
    total_issues: 0,
    trigger: ""
  },
  analyzes: []
}

const analysisSlice = createSlice({
  name: 'analysis',
  initialState: initialState as AnalysisState,
  reducers: {
    started: (state, { payload }) => {
      const data = {
        id: payload.id,
        created_at: payload.created_at,
        created_by: payload.created_by,
        duration: payload.duration,
        finished_at: payload.finished_at,
        last_commit_hash: payload.last_commit_hash,
        repository: payload.repository,
        severity: payload.severity,
        started_at: payload.started_at,
        state: payload.state,
        total_issues: payload.total_issues,
        trigger: payload.trigger
      }
      let index = state.analyzes.findIndex(analysis => analysis.id === payload.id)
      if (index !== -1) {
        state.analyzes[index] = data
      }
    }
  },
  extraReducers: (builder) => {
    builder
      .addCase(getAnalysisById.fulfilled, (state, { payload }) => {
        state.analysis = payload
      })
      .addCase(getAnalyzes.fulfilled, (state, { payload }) => {
        state.analyzes = payload
      })
  }
})

const { reducer } = analysisSlice
export default reducer
