import { createSlice } from '@reduxjs/toolkit'

export interface NotificationState {
  isEstablishingConnection: boolean
  isConnected: boolean
}

const initialState: NotificationState = {
  isEstablishingConnection: false,
  isConnected: false
}

const notificationSlice = createSlice({
  name: 'notification',
  initialState,
  reducers: {
    startConnecting: (state => {
      state.isEstablishingConnection = true
    }),
    connectionEstablished: (state => {
      state.isConnected = true
      state.isEstablishingConnection = true
    })
  }
})

export const { startConnecting, connectionEstablished } = notificationSlice.actions
const { reducer } = notificationSlice
export default reducer
