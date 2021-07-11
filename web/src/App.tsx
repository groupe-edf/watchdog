import { createBrowserHistory } from 'history'
import { Component } from 'react';
import { Provider } from 'react-redux';
import {
  BrowserRouter,
  Switch,
  Route
} from "react-router-dom";
import {
  Login,
  NotFound,
} from './pages';
import { Layout } from './components/Layout';
import { SOCKET_URL } from './constants';
import configureStore from './configureStore';
import { routes } from './routes';

const history = createBrowserHistory()
const initialState = window.INITIAL_REDUX_STATE
const store = configureStore(history, initialState)

class App extends Component {
  constructor(props: any) {
    super(props);
  }
  componentDidMount() {
    const ws = new WebSocket(SOCKET_URL);
    ws.onmessage = (evt: MessageEvent) => {
    };
  }
  render() {
    return (
      <Provider store={store}>
        <BrowserRouter>
          <Switch>
            {routes.map(({component: Component, exact, isPrivate, path}, key) => (
              <Route
                exact={exact}
                key={key}
                path={path}
                render={props => isPrivate ? <Layout title="Welcome to Watchdog"><Component {...props} /></Layout> : <Component {...props} />}></Route>
            ))}
          </Switch>
        </BrowserRouter>
      </Provider>
    );
  }
}

export default App;
