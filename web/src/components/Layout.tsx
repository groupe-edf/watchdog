import { ReactNode } from "react"
import {
  Box,
  Container,
  Text,
  useColorModeValue,
} from '@chakra-ui/react';
import Sidebar from "./Sidebar"

const layouts = {
}

type Props = {
  children?: ReactNode
  title?: string
}

export const Layout = ({
  children,
  title = "This is the default title",
}: Props) => (
  <Sidebar>
    {children}
    <Box
      borderTopWidth={1}
      borderStyle={'solid'}
      borderColor={useColorModeValue('gray.200', 'gray.700')}
      mt={4}>
      <Container
        maxW={'6xl'}
        py={4}
        direction={{ base: 'column', md: 'row' }}
        spacing={4}
        justify={{ md: 'space-between' }}
        align={{ md: 'center' }}>
          <Text align="center">Watchdog</Text>
      </Container>
    </Box>
  </Sidebar>
)
