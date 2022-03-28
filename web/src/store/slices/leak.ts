import { createAsyncThunk, createSlice } from "@reduxjs/toolkit"
import { Leak, Query } from "../../models"
import { LeakService } from "../../services"

export enum LeakActionTypes {
  LEAKS_FIND = "@@leaks/FIND",
  LEAKS_FIND_BY_ID = "@@leaks/FIND_BY_ID",
}

export interface LeakState {
  readonly leak: Leak
  readonly leaks: Leak[]
}

export const getLeaks = createAsyncThunk(
  LeakActionTypes.LEAKS_FIND,
  async (query: Query, thunkAPI) => {
    const response = await LeakService.findAll(query)
    return response.data
  }
)

export const getLeakById = createAsyncThunk(
  LeakActionTypes.LEAKS_FIND_BY_ID,
  async (leakId: string, thunkAPI) => {
    try {
      const response = await LeakService.findById(leakId)
      return response.data
    } catch (error) {
      return thunkAPI.rejectWithValue(error)
    }
  }
)

export const initialState: LeakState = {
  leak: {
    id: "",
    author_email: "",
    author_name: "",
    commit_hash: "",
    created_at: "",
    file: "",
    line: "",
    line_number: 0,
    occurence: 0,
    offender: "",
    repository: {
      id: "",
      repository_url: "",
      visibility: ""
    },
    rule: {
      id: 0,
      display_name: "",
      enabled: true,
      severity: "",
      tags: [],
    },
    severity: ""
  },
  leaks: []
}

const leakSlice = createSlice({
  name: 'issue',
  initialState: initialState as LeakState,
  reducers: {
  },
  extraReducers: (builder) => {
    builder
      .addCase(getLeaks.fulfilled, (state, { payload }) => {
        state.leaks = payload
      })
      .addCase(getLeakById.fulfilled, (state, { payload }) => {
        state.leak = payload
      })
  }
})

const { reducer } = leakSlice
export default reducer
