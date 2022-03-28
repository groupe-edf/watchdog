import { Fragment, StrictMode } from "react"
import ReactDOM from "react-dom/client"
import App from "./App"
import {
  ChakraProvider,
  ColorModeScript,
  extendTheme
} from "@chakra-ui/react"
import store from './configureStore'
import { Provider } from 'react-redux'
import reportWebVitals from './reportWebVitals'
import interceptor from './common/interceptor'
import { getCategories, getSettings } from "./store/slices/global"
import { BrowserRouter } from "react-router-dom"
import { createBrowserHistory } from 'history'
import { startConnecting } from "./store/slices/notification"

import { SaasProvider, theme as baseTheme } from '@saas-ui/react'
import extendedTheme from "./theme"
const theme = extendTheme(baseTheme, extendedTheme)

const history = createBrowserHistory()
interceptor.interceptor(store, history)
store.dispatch(getCategories())
store.dispatch(getSettings())
store.dispatch(startConnecting())

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
)
root.render(
  <Fragment>
    <ColorModeScript initialColorMode={extendedTheme.config.initialColorMode} />
    <SaasProvider theme={theme} resetCSS={true}>
      <Provider store={store}>
        <BrowserRouter>
          <App />
        </BrowserRouter>
      </Provider>
    </SaasProvider>
  </Fragment>
)
reportWebVitals()
