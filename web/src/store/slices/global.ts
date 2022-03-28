import Axios from 'axios'
import { createAsyncThunk, createSlice } from "@reduxjs/toolkit"
import { CategoryService, GlobalService } from '../../services'
import { Category } from '../../models'
import { Pagination } from '../../models/common'

export enum GlobalActionTypes {
  GLOBAL_GET_CATEGORIES = "@@global/GET_CATEGORIES",
  GLOBAL_GET_MESSAGE = "@@global/GET_MESSAGE",
  GLOBAL_GET_SETTINGS = "@@global/GET_SETTINGS",
  GLOBAL_GET_VERSION = "@@global/GET_VERSION",
  GLOBAL_SET_MESSAGE = "@@global/SET_MESSAGE"
}

export interface Global {
  api_keys: any
  categories: Category[],
  current_api_key: any
  loading: boolean
  message: any
  pagination: Pagination
  settings: any
  version: any
}
const initialState: Global = {
  api_keys: [],
  categories: [],
  current_api_key: {},
  loading: false,
  message: "",
  pagination: {
    currentPage: 1,
    pagesToShow: 5,
    itemsPerPage: 10,
    offset: 0,
    totalItems: 0
  },
  version: {},
  settings: {}
}
export const getCategories = createAsyncThunk(
  GlobalActionTypes.GLOBAL_GET_CATEGORIES,
  async () => {
    const response = await CategoryService.findAll()
    return response.data
  }
)

export const getSettings = createAsyncThunk(
  GlobalActionTypes.GLOBAL_GET_SETTINGS,
  async () => {
    const response = await Axios.get('/settings');
    return response.data
  }
)

export const getVersion = createAsyncThunk(
  GlobalActionTypes.GLOBAL_GET_VERSION,
  async (any, thunkAPI) => {
    try {
      const response = await GlobalService.getVersion()
      return response.data
    } catch (error) {
      return thunkAPI.rejectWithValue(error)
    }
  }
)

export const globaleSlice = createSlice({
  name: 'global',
  initialState: initialState as Global,
  reducers: {
    setPagination: (state, { payload }) => {
      state.pagination = payload
    },
    setLoading: (state, { payload }) => {
      state.loading = payload
    },
    setMessage: (state, { payload }) => {
      state.message = payload
    },
    clearMessage: (state) => {
      state.message = ""
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(getCategories.fulfilled, (state, { payload }) => {
        state.categories = payload
      })
      .addCase(getSettings.fulfilled, (state, { payload }) => {
        state.settings = payload
      })
      .addCase(getVersion.fulfilled, (state, { payload }) => {
        state.version = payload
      })
  }
})

const { reducer, actions } = globaleSlice;
export const { clearMessage, setLoading, setMessage, setPagination } = actions
export default reducer
