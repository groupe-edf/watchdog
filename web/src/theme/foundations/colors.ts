import { theme } from "@chakra-ui/react"

const brand = {
  100: "#f14e32",
  200: "#f14e32",
  300: "#f14e32",
  400: "#f14e32",
  500: "#f14e32",
  600: "#f14e32",
  700: "#f14e32",
  800: "#f14e32",
  900: "#f14e32"
}
const font = {
  primary: {
    lightMode: theme.colors.gray["700"],
    darkMode: theme.colors.gray["200"],
  },
  secondary: {
    lightMode: theme.colors.gray["600"],
    darkMode: theme.colors.gray["400"],
  }
}

const colors = {
  brand,
  font
}

export default colors
