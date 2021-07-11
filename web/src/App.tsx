import { Component } from 'react';
import { connect, ConnectedProps } from 'react-redux';
import {
  BrowserRouter,
  Switch,
  Redirect
} from "react-router-dom";
import { API_PATH, SOCKET_URL } from './constants';
import { routes } from './routes';
import { PrivateRoute } from './components/PrivateRoute';
import { ApplicationState } from './store';
import { CategoryActionTypes } from './store/categories/types';
import categoryService from "./services/category"
import { fetchData } from './services/commons';
import { GlobalActionTypes } from './store/global/types';

const mapState = (state: ApplicationState) => ({
  state: state.repositories
})
const mapDispatch = {
  getSettings: (payload: any) => ({ type: GlobalActionTypes.GLOBAL_SETTINGS, payload }),
  getCategories: (payload: any) => ({ type: CategoryActionTypes.CATEGORY_FIND_ALL, payload })
}
const connector = connect(mapState, mapDispatch)
type ApplicationProps = ConnectedProps<typeof connector>

class App extends Component<ApplicationProps> {
  constructor(props: ApplicationProps) {
    super(props)
  }
  componentDidMount() {
    const ws = new WebSocket(SOCKET_URL)
    ws.onmessage = (event: MessageEvent) => {}
    const { getCategories, getSettings } = this.props
    categoryService.findAll().then(response => {
      getCategories(response.data)
    })
    fetchData("GET", `${API_PATH}/settings`).then(response => {
      getSettings(response.data)
    })
  }
  render() {
    return (
      <BrowserRouter>
        <Switch>
          {routes.map(({component: Component, exact, isPrivate, path}, key) => (
            <PrivateRoute
              exact={exact}
              key={key}
              path={path}
              isPrivate={isPrivate}
              component={Component}></PrivateRoute>
          ))}
          <Redirect from="*" to="/" />
        </Switch>
      </BrowserRouter>
    );
  }
}

export default connector(App)
