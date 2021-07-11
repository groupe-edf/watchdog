import { Reducer } from 'redux';
import { WebsocketState, WebsocketActions, WebsocketActionTypes } from './types';

export const initialState: WebsocketState = {
  messages: []
}

const reducer: Reducer<WebsocketState> = (state: WebsocketState = initialState, action) => {
  const { type, payload } = action;
  switch ((action as WebsocketActions).type) {
    case WebsocketActionTypes.WEBSOCKET_CONNECT:
      return state;
    default:
      return state;
  }
}

export default reducer;
