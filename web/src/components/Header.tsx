import {
  Flex,
  Menu,
  MenuButton,
  HStack,
  Avatar,
  MenuList,
  MenuItem,
  MenuDivider,
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  Text,
  Link,
  SkeletonCircle,
  useColorMode,
  useColorModeValue,
  Button} from "@chakra-ui/react"
import { IoChevronForwardOutline, IoLogoGithub, IoLogOutOutline, IoMoonOutline, IoPersonOutline, IoSettingsOutline, IoSunnyOutline } from "react-icons/io5"
import { useSelector } from "react-redux"
import { Link as RouterLink } from "react-router-dom"
import { RootState } from "../configureStore"
import authenticationService from "../services/authentication"

const Header = () => {
  const { colorMode, toggleColorMode } = useColorMode()
  const { currentUser } = useSelector((state: RootState) => state.authentication)
  return (
    <Flex
      alignItems="center"
      as="header"
      background={useColorModeValue('white', 'gray.900')}
      borderBottomWidth="1px"
      borderBottomColor={useColorModeValue('gray.200', 'gray.700')}
      color={useColorModeValue('gray.800', 'white')}
      justify="space-between"
      paddingX={{ base: 4, md: 4 }}
      paddingY={{ base: 2, md: 2 }}>
      <Breadcrumb separator={<IoChevronForwardOutline color="gray.500" />}>
        <BreadcrumbItem isCurrentPage>
          <BreadcrumbLink as={Link} to="#" href="#">Home</BreadcrumbLink>
        </BreadcrumbItem>
      </Breadcrumb>
      <Button onClick={toggleColorMode}>
        {colorMode === 'light' ? <IoMoonOutline/> : <IoSunnyOutline/>}
      </Button>
      <Flex align="center">
        {currentUser ? (
        <Menu>
          <MenuButton>
            <HStack>
              <Text>Welcome</Text>
              <Text fontWeight="bold">{currentUser.last_name}</Text>
              <Avatar background="brand.100" name={currentUser.first_name} size="sm" />
            </HStack>
          </MenuButton>
          <MenuList
            background="white"
            borderColor="gray.200">
            <MenuItem icon={<IoPersonOutline/>} as={RouterLink} to="/profile">
              Profile
            </MenuItem>
            <MenuItem icon={<IoSettingsOutline/>} as={RouterLink} to="/settings">
              Settings
            </MenuItem>
            <MenuDivider/>
            <MenuItem icon={<IoLogOutOutline/>} as={RouterLink} to="/login" color="red" onClick={authenticationService.logout}>
              Logout
            </MenuItem>
          </MenuList>
        </Menu>
        ) : (
          <SkeletonCircle size="10" />
        )}
      </Flex>
    </Flex>
  )
}

export default Header
