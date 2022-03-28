import { createAsyncThunk, createSlice } from "@reduxjs/toolkit"
import { AxiosError } from "axios"
import { Analysis, Query, Repository } from "../../models"
import { AnalysisService, RepositoryService } from '../../services'

export enum RepositoryActionTypes {
  REPOSITORIES_FIND = "@@repositories/FIND",
  REPOSITORIES_FIND_BY_ID = "@@repositories/FIND_BY_ID",
  REPOSITORIES_ANALYZE = "@@repositories/ANALYZE"
}

export const analyze = createAsyncThunk(
  RepositoryActionTypes.REPOSITORIES_ANALYZE,
  async (data: { repository_url: string }, thunkAPI) => {
    try {
      const response = await RepositoryService.analyze(data)
      thunkAPI.dispatch(getRepositories({}))
      return response.data
    } catch (error) {
      return thunkAPI.rejectWithValue(error)
    }
  }
)

export const analyzeRepository = createAsyncThunk(
  RepositoryActionTypes.REPOSITORIES_ANALYZE,
  async (data: {repository_id: string}, thunkAPI) => {
    try {
      const response = await RepositoryService.analyzeById(data.repository_id, {})
      return response.data
    } catch (error) {
      if (error instanceof AxiosError) {
        const message = ( error.response && error.response.data && error.response.data.detail ) || error.message || error.toString()
        return thunkAPI.rejectWithValue(message)
      }
      throw error
    }
  }
)

export const getRepositories = createAsyncThunk(
  RepositoryActionTypes.REPOSITORIES_FIND,
  async (query: Query, thunkAPI) => {
    try {
      const response = await RepositoryService.findAll(query)
      return response.data
    } catch (error) {
      return thunkAPI.rejectWithValue(error)
    }
  }
)

export const getRepositoryById = createAsyncThunk(
  RepositoryActionTypes.REPOSITORIES_FIND_BY_ID,
  async (repositoryId: string, thunkAPI) => {
    try {
      const response = await RepositoryService.findById(repositoryId)
      return response.data
    } catch (error) {
      return thunkAPI.rejectWithValue(error)
    }
  }
)

export interface RepositoryState {
  readonly repositories: Repository[]
  readonly repository: Repository
}

export const initialState: RepositoryState = {
  repositories: [],
  repository: {
    id: "",
    repository_url: "",
    visibility: ""
  }
}

const repositorySlice = createSlice({
  name: 'repository',
  initialState: initialState as RepositoryState,
  reducers: {
    setAnalysis: (state, { payload }) => {
      if (state.repository.id === payload.repository.id) {
        state.repository.last_analysis = payload
      }
      let index = state.repositories.findIndex(repository => repository.id === payload.repository.id)
      if (index !== -1) {
        state.repositories[index].last_analysis = payload
      }
    }
  },
  extraReducers: (builder) => {
    builder
      .addCase(analyzeRepository.fulfilled, (state, { payload }) => {
        if (state.repository.id === payload.repository?.id) {
          state.repository.last_analysis = payload
        }
        const index = state.repositories.findIndex(repository => repository.id === payload.repository?.id)
        if (index !== -1) {
          state.repositories[index].last_analysis = payload
        }
      })
      .addCase(getRepositories.fulfilled, (state, { payload }) => {
        state.repositories = payload
      })
      .addCase(getRepositoryById.fulfilled, (state, { payload }) => {
        state.repository = payload
      })
  }
})

const { reducer } = repositorySlice
export default reducer
