import { SystemStyleFunction, mode } from "@chakra-ui/theme-tools"
import colors from "../foundations/colors"

const baseStyle = (props: SystemStyleFunction) => {
  return {
    color: mode(colors.font.primary.lightMode, colors.font.primary.darkMode)(props)
  }
}

export default {
  baseStyle
}
