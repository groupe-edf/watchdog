import { Badge, Tooltip } from "@chakra-ui/react"

export const StatusBadge = ({ state, hint }: { state: string, hint?: string }) => {
  let colorScheme = "secondary"
  let text = ""
  switch (state) {
    case "failed":
      colorScheme = "red"
      text = "Failed"
      break
    case "started":
      colorScheme = "orange"
      text = "Started"
      break
    case "success":
      colorScheme = "green"
      text = "Success"
      break
  }
  return (
    <Tooltip label={hint}>
      <Badge colorScheme={colorScheme}>{text}</Badge>
    </Tooltip>
  )
}
