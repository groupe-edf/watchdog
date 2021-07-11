import { createBrowserHistory } from 'history'
import { StrictMode } from "react"
import { render } from "react-dom"
import App from "./App"
import {
  ChakraProvider,
  ColorModeScript,
  extendTheme
} from "@chakra-ui/react"
import configureStore from './configureStore';
import { Provider } from 'react-redux'

const history = createBrowserHistory()
const initialState = window.INITIAL_REDUX_STATE
const store = configureStore(history, initialState)

const theme = extendTheme({
  colors: {
    brand: {
      100: "#f14e32",
      200: "#f14e32",
      300: "#f14e32",
      400: "#f14e32",
      500: "#f14e32",
      600: "#f14e32",
      700: "#f14e32",
      800: "#f14e32",
      900: "#f14e32",
    },
  },
  components: {
    baseStyle: {},
    Button: {
    },
    Input: {
      baseStyle: {
        focusBorderColor: "brand.100"
      }
    },
    Select: {
      baseStyle: {
        focusBorderColor: "brand.100"
      }
    }
  },
  modifiers: {
    computeStyle: {
      gpuAcceleration: false
    }
  },
  styles: {
    global: (props) => ({
      "html, body": {
        overflowY: "hidden"
      }
    })
  }
})

render(
  <StrictMode>
    <ColorModeScript initialColorMode={'dark'} />
    <ChakraProvider theme={theme}>
      <Provider store={store}>
        <App />
      </Provider>
    </ChakraProvider>
  </StrictMode>,
  document.getElementById("root"),
)
