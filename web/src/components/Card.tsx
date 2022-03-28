import { Box, useColorModeValue, useStyleConfig } from '@chakra-ui/react'

export const Card = (props: any) => {
  const { variant, children, ...rest } = props
  const styles = useStyleConfig("Card", { variant })
  return (
    <Box __css={styles} {...rest}>{children}</Box>
  )
}
