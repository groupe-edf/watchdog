import {
  Store,
  applyMiddleware,
  createStore,
} from 'redux'
import { composeWithDevTools } from 'redux-devtools-extension'
import { History } from 'history';
import { routerMiddleware } from 'connected-react-router'
import { ApplicationState, createRootReducer } from './store';

const store = function configureStore(history: History, initialState: ApplicationState): Store<ApplicationState> {
  const store = createStore(
    createRootReducer(history),
    initialState as any,
    composeWithDevTools(
      applyMiddleware(
        routerMiddleware(history)
      )
    )
  );
  return store
}

export default store;
