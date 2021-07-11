import { Component, useState } from "react"
import {
  IoSettingsOutline
} from "react-icons/io5"
import {
  Box,
  BoxProps,
  Button,
  Divider,
  Flex,
  Icon,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
  Link,
  Text,
  useColorModeValue,
  useDisclosure
} from "@chakra-ui/react"
import { ArrowUpDownIcon } from "@chakra-ui/icons"
import { connect } from "react-redux"
import { NavLink } from "react-router-dom"
import { routes } from "../../routes"
import Header from "../Header"
import { ItemContent, NavItem } from "./NavItem"
import CurrentUser from "./CurrentUser"

interface SidebarContentProps extends BoxProps {
  onClose?: () => void;
}

const SidebarContent = ({ onClose, ...rest }: SidebarContentProps) => {
  const settings = useDisclosure()
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
        <Text align="center" width="100%" color="brand.100" fontSize="2xl" fontFamily="monospace">
          Watchdog
        </Text>
        <Divider marginY={4} display={size == "small" ? "none" : "flex"} />
        {routes.map((route) => (
          !route.hide && <NavItem key={route.title} route={route.path} icon={route.icon} size={size}>
            {route.title}
          </NavItem>
        ))}
        <Menu placement="right">
          <Link
            borderRadius={8}
            _hover={{ background: "brand.100", color: "white" }}
            style={{ textDecoration: "none" }}
            width="100%">
            <MenuButton width="100%">
              <ItemContent icon={IoSettingsOutline} route="" children="Settings" size="large"></ItemContent>
            </MenuButton>
          </Link>
          <MenuList>
            <MenuItem as={NavLink} to="/integrations">Integrations</MenuItem>
            <MenuItem as={NavLink} to="/jobs">Jobs</MenuItem>
            <MenuItem as={NavLink} to="/users">Users</MenuItem>
            <MenuItem as={NavLink} to="/workspaces">Workspaces</MenuItem>
          </MenuList>
        </Menu>
      </Flex>
      <Flex
        flexDir="column"
        width="100%"
        alignItems={size == "small" ? "center" : "flex-start"}>
        <Menu matchWidth placement="right">
          <MenuButton as={Button} colorScheme="brand" width="100%">
            <Flex justifyContent="space-between">
              <Text>Workspace</Text>
              <Icon as={ArrowUpDownIcon} />
            </Flex>
          </MenuButton>
          <MenuList>
            <MenuItem>Default</MenuItem>
          </MenuList>
        </Menu>
        <Divider marginY={4} display={size == "small" ? "none" : "flex"} />
        <CurrentUser />
      </Flex>
    </Flex>
  )
}

export class Sidebar extends Component<any> {
  constructor(props: any) {
    super(props);
  }
  render() {
    return (
      <Flex
        width="100%"
        background="gray.100"
        height="100vh">
        <SidebarContent display={{ base: 'none', md: 'block' }}/>
        <Flex
          flexDirection="column"
          width="100%"
          overflowY="auto">
          <Header/>
          <Box padding={4}>{this.props.children}</Box>
        </Flex>
      </Flex>
    )
  }
}

export default connect()(Sidebar);
