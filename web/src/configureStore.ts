import {
  Store,
  applyMiddleware,
  compose,
  createStore,
} from 'redux'
import { History } from 'history';
import { composeWithDevTools } from 'redux-devtools-extension';
import { ApplicationState, reducers } from './store';

const store = function configureStore(
  history: History,
  initialState: ApplicationState,
): Store<ApplicationState> {
  const composeEnhancers = composeWithDevTools({});
  const store = createStore(
    reducers,
    initialState,
    compose(applyMiddleware())
  );
  return store
}

export default store;
