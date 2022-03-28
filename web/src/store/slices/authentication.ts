import { createSlice, createAsyncThunk } from '@reduxjs/toolkit'
import { setMessage } from './global'
import AuthenticationService from '../../services/authentication'

export enum AuthenticationActionTypes {
  AUTHENTICATION_CURRENT_USER = "@@authentication/CURRENT_USER",
  AUTHENTICATION_GET_PROFILE = "@@authentication/GET_PROFILE",
  AUTHENTICATION_LOGIN = "@@authentication/LOGIN",
  AUTHENTICATION_LOGOUT = "@@authentication/LOGOUT",
  AUTHENTICATION_REGISTER = "@@authentication/REGISTER"
}

export const getProfile = createAsyncThunk(
  AuthenticationActionTypes.AUTHENTICATION_GET_PROFILE,
  async () => {
    const response = await AuthenticationService.getProfile()
    localStorage.setItem("user", JSON.stringify(response.data))
    return response.data
  }
)

export const login = createAsyncThunk(
  AuthenticationActionTypes.AUTHENTICATION_LOGIN,
  async ({ email, password }: any, thunkAPI) => {
    try {
      const response = await AuthenticationService.login({ email, password })
      if (response.data.token) {
        localStorage.setItem("token", response.data.token)
      }
      return response.data
    } catch (error: any) {
      return thunkAPI.rejectWithValue(error?.response?.data)
    }
  }
)

export const logout = createAsyncThunk(
  AuthenticationActionTypes.AUTHENTICATION_LOGOUT,
  async () => {
    AuthenticationService.logout()
  }
)
export const register = createAsyncThunk(
  AuthenticationActionTypes.AUTHENTICATION_REGISTER,
  async (values: any, thunkAPI) => {
    try {
      const response = await AuthenticationService.register(values)
      thunkAPI.dispatch(setMessage(response.data.message))
      return response.data
    } catch (error: any) {
      return thunkAPI.rejectWithValue(error?.response?.data || error)
    }
  }
)
export interface Authentication {
  accessToken?: string
  currentUser?: any
  isLoggedIn: boolean
}
const initialState: Authentication = { isLoggedIn: false }

const authenticationSlice = createSlice({
  name: 'authentication',
  initialState: initialState as Authentication,
  reducers: {
    logout: () => initialState
  },
  extraReducers: (builder) => {
    builder
      .addCase(getProfile.fulfilled, (state, { payload }) => {
        state.currentUser = payload
      })
      .addCase(login.fulfilled, (state, { payload }) => {
        state.accessToken = payload.token
        state.isLoggedIn = true
      })
      .addCase(login.rejected, (state, { payload }) => {
        state.isLoggedIn = false
      })
      .addCase(logout.fulfilled, (state, { payload }) => {
        state = initialState
      })
  }
})

const { reducer } = authenticationSlice
export default reducer
