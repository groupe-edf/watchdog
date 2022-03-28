import { useState } from "react"
import {
  Box,
  BoxProps,
  Divider,
  Flex,
  Text,
  useColorModeValue
} from "@chakra-ui/react"
import { routes } from "../../routes"
import { NavItem } from "./NavItem"
import Header from "../Header"
import { Link } from "react-router-dom"

interface SidebarContentProps extends BoxProps {
  onClose?: () => void
}

const SidebarContent = ({ onClose, ...rest }: SidebarContentProps) => {
  const [size, setSize] = useState("large")
  return (
    <Flex
      as="aside"
      background={useColorModeValue('white', 'gray.900')}
      borderRight="1px"
      borderRightColor={useColorModeValue('gray.200', 'gray.700')}
      position="sticky"
      zIndex="sticky"
      padding={2}
      width={size == "small" ? "75px" : "200px"}
      flexDir="column"
      justifyContent="space-between">
      <Flex
        flexDir="column"
        width="100%"
        alignItems={size == "small" ? "center" : "flex-start"}
        as="nav">
        <Text as={Link} to="/" align="center" width="100%" color="brand.100" fontSize="2xl" fontFamily="monospace">
          Watchdog
        </Text>
        <Divider marginY={4} display={size == "small" ? "none" : "flex"} />
        {routes.map((route) => (
          !route.hide && <NavItem key={route.path} route={route.path} icon={route.icon} size={size}>
            {route.title}
          </NavItem>
        ))}
      </Flex>
    </Flex>
  )
}

const Sidebar = (props: any) => {
  const { children } = props
  return (
    <Flex
      background={useColorModeValue('gray.100', 'gray.900')}
      width="100%"
      height="100vh">
      <SidebarContent display={{ base: 'none', md: 'block' }}/>
      <Flex
        flexDirection="column"
        width="100%"
        overflowY="auto">
        <Header/>
        <Box padding={4}>{children}</Box>
      </Flex>
    </Flex>
  )
}

export default Sidebar
