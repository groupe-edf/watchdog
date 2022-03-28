import { SystemStyleFunction, mode } from "@chakra-ui/theme-tools"

const baseStyle = {
  borderRadius: "lg",
  borderWidth: "1px",
  padding: 6
}
const variants = {
  panel: (props: any) => ({
    background: mode("white", "gray.800")(props),
  })
}
const defaultProps = {
  variant: "panel"
}

export default {
  baseStyle,
  defaultProps,
  variants
}
