import { ActionCreator } from 'redux';
import { WebsocketActionTypes, WebsocketConnectAction, WebsocketState } from './types';

export const connect: ActionCreator<WebsocketConnectAction> = (state: WebsocketState) => ({
  type: WebsocketActionTypes.WEBSOCKET_CONNECT,
  payload: {
    state,
  },
});
