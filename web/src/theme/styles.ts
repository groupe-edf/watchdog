import { Styles } from "@chakra-ui/theme-tools"
import { mode } from "@chakra-ui/theme-tools"
import colors from "./foundations/colors"

const styles: Styles = {
  global: (props: any) => ({
    body: {
      color: mode(colors.font.primary.lightMode, colors.font.primary.darkMode)(props)
    }
  })
}

export default styles
