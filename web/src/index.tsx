import { StrictMode } from "react"
import { render } from "react-dom"
import App from "./App"
import {
  ChakraProvider,
  ColorModeScript,
  extendTheme
} from "@chakra-ui/react"

const theme = extendTheme({
  colors: {
    git: {
      background: "#f0efe7",
      default: "#4e443c",
      link: "#0388a6",
      primary: "#f14e32"
    }
  },
  components: {
    Button: {
    }
  },
  styles: {
  }
})

render(
  <StrictMode>
    <ColorModeScript initialColorMode={'dark'} />
    <ChakraProvider theme={theme}>
      <App />
    </ChakraProvider>
  </StrictMode>,
  document.getElementById("root"),
)
