import { Badge, Tooltip } from "@chakra-ui/react"

export const StatusBadge = ({ state, hint, ...rest }: { state: string, hint?: string, rest?: any }) => {
  let colorScheme = "gray"
  let text = state.toUpperCase()
  switch (state) {
    case "failed":
      colorScheme = "red"
      break
    case "in_progress":
      colorScheme = "yellow"
      break
    case "started":
      colorScheme = "orange"
      break
    case "success":
      colorScheme = "green"
      break
  }
  return (
    <Tooltip label={hint}>
      <Badge colorScheme={colorScheme} {...rest}>{text}</Badge>
    </Tooltip>
  )
}
