import { Middleware } from "@reduxjs/toolkit"
import { SOCKET_URL } from "../../constants"
import { connectionEstablished, startConnecting } from "../slices/notification"

const notificationMiddleware: Middleware = store => {
  let socket: WebSocket
  return next => action => {
    const state = store.getState().notifications
    const isConnectionEstablished = socket && store.getState().notifications.isConnected
    if (action.type === startConnecting.type) {
      socket = new WebSocket(SOCKET_URL)
      socket.onopen = (event) => {
        store.dispatch(connectionEstablished())
      }
      socket.onmessage = (event) => {
        const data = JSON.parse(event.data)
        console.log(data.event_type)
        store.dispatch({
          type: data.event_type.replace(':', '/'),
          payload: data.payload
        })
      }
    }
    next(action)
  }
}
export default notificationMiddleware
