import { Component } from 'react';
import { Route, Redirect } from 'react-router-dom';
import { Layout } from './Layout';

class PrivateRoute extends Component<any> {
  render() {
    const { component: Component, isPrivate, ...rest } = this.props
    return (
      <Route {...rest} render={props => {
        if (isPrivate) {
          if (!localStorage.getItem('user')) {
            return <Redirect to={{ pathname: '/login', state: { from: props.location } }} />
          }
          return <Layout><Component {...props} /></Layout>
        }
        return <Component {...props} />
      }} />
    )
  }
}

export { PrivateRoute };
