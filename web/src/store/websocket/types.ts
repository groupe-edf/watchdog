import { Action } from "redux";

export enum WebsocketActionTypes {
  WEBSOCKET_BEGIN_RECONNECT = "@@websocket/BEGIN_RECONNECT",
  WEBSOCKET_BROKEN = "@@websocket/BROKEN",
  WEBSOCKET_CLOSED = "@@websocket/CLOSED",
  WEBSOCKET_CONNECT = "@@websocket/CONNECT",
  WEBSOCKET_DISCONNECT = "@@websocket/DISCONNECT",
  WEBSOCKET_ERROR = "@@websocket/ERROR",
  WEBSOCKET_MESSAGE = "@@websocket/MESSAGE",
  WEBSOCKET_RECONNECT_ATTEMPT = "@@websocket/RECONNECT_ATTEMPT",
  WEBSOCKET_RECONNECTED = "@@websocket/RECONNECTED",
  WEBSOCKET_SEND = "@@websocket/SEND",
}

export interface WebsocketState {
  messages: string[];
}

export interface WebsocketConnectAction extends Action {
  type: WebsocketActionTypes.WEBSOCKET_CONNECT;
  payload: {
    state: WebsocketState;
  };
}

export interface WebsocketMessageAction extends Action {
  type: WebsocketActionTypes.WEBSOCKET_MESSAGE;
  payload: {
    state: WebsocketState;
  };
}

export type WebsocketActions = WebsocketConnectAction | WebsocketMessageAction;
