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
  Icon
} from "@chakra-ui/react"
import { Component } from "react"
import { IoChevronForwardOutline, IoLogoGithub, IoLogOutOutline, IoPersonOutline, IoSettingsOutline } from "react-icons/io5"
import { connect, ConnectedProps } from "react-redux"
import { Link as RouterLink, RouteComponentProps } from "react-router-dom"
import { API_PATH } from "../constants"
import authenticationService from "../services/authentication"
import { fetchData } from "../services/commons"
import { ApplicationState } from "../store"
import { AuthenticationActionTypes } from "../store/authentication/types"

const mapState = (state: ApplicationState) => ({
  currentUser: state.authentication.currentUser
})
const mapDispatch = {
  getCurrentUser: (payload: any) => ({ type: AuthenticationActionTypes.AUTHENTICATION_CURRENT_USER, payload }),
}
const connector = connect(mapState, mapDispatch)
type HeaderProps = ConnectedProps<typeof connector> & RouteComponentProps

export class Header extends Component<HeaderProps, {
  isLoading: boolean
}> {
  constructor(props: HeaderProps) {
    super(props);
    this.state = {
      isLoading: false
    }
  }
  componentDidMount() {
    const { currentUser, getCurrentUser } = this.props
    if (currentUser.email === "") {
      this.setState({ isLoading: true })
      fetchData("GET", `${API_PATH}/profile`)
        .then(response => {
          getCurrentUser(response.data)
          this.setState({ isLoading: false })
        })
        .catch(response => {
          authenticationService.logout()
        })
    }
  }
  render() {
    const { currentUser } = this.props;
    return (
      <Flex
        as="header"
        paddingX={{ base: 4, md: 4 }}
        paddingY={{ base: 2, md: 2 }}
        alignItems="center"
        background="white"
        borderBottomWidth="1px"
        borderBottomColor="gray.200"
        justify="space-between">
        <Breadcrumb separator={<IoChevronForwardOutline color="gray.500" />}>
          <BreadcrumbItem isCurrentPage>
            <BreadcrumbLink as={Link} to="#" href="#">Home</BreadcrumbLink>
          </BreadcrumbItem>
        </Breadcrumb>
        <Flex align="center">
          {!this.state.isLoading ? (
          <Menu>
            <MenuButton>
              <HStack>
                <Text>Bonjour</Text>
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
              <MenuItem icon={<IoLogoGithub/>} as={Link} href="https://github.com/groupe-edf/watchdog" isExternal>
                groupe-edf/watchdog
              </MenuItem>
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
}

export default connector(Header as any)
