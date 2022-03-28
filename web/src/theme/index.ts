import { extendTheme, type ThemeConfig, ThemeDirection } from "@chakra-ui/react"
import styles from "./styles"
import components from "./components"
import * as foundations from "./foundations"

const direction: ThemeDirection = "ltr"
const config: ThemeConfig = {
  initialColorMode: 'light',
  useSystemColorMode: false,
  cssVarPrefix: "watchdog"
}

const themeSettings  = {
  components,
  direction,
  ...foundations,
  ...styles,
  config
}

export default extendTheme(themeSettings)
