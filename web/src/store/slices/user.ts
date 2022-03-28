import { createAsyncThunk, createSlice } from "@reduxjs/toolkit"
import { User } from "../../models"
import { UserService } from "../../services"

export enum UserActionTypes {
  USERS_FIND = "@@users/FIND",
  USERS_FIND_BY_ID = "@@users/FIND_BY_ID"
}

export const getUsers = createAsyncThunk(
  UserActionTypes.USERS_FIND,
  async () => {
    const response = await UserService.findAll();
    return response.data
  }
)

export interface UserState {
  readonly users: User[]
}

export const initialState: UserState = {
  users: [],
}

const repositorySlice = createSlice({
  name: 'user',
  initialState: initialState as UserState,
  reducers: {
  },
  extraReducers: (builder) => {
    builder
      .addCase(getUsers.fulfilled, (state, { payload }) => {
        state.users = payload
      })
  }
})

const { reducer } = repositorySlice
export default reducer
