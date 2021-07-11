import { Avatar, Flex, Heading, Text } from "@chakra-ui/react"
import { Component } from "react"
import { connect, ConnectedProps } from "react-redux"
import { RouteComponentProps } from "react-router"
import { API_PATH } from "../../constants"
import { fetchData } from "../../services/commons"
import { ApplicationState } from "../../store"
import { AuthenticationActionTypes } from "../../store/authentication/types"
import authenticationService from "../../services/authentication"

const mapState = (state: ApplicationState) => ({
  currentUser: state.authentication.currentUser
})
const mapDispatch = {
  getCurrentUser: (payload: any) => ({ type: AuthenticationActionTypes.AUTHENTICATION_CURRENT_USER, payload }),
}
const connector = connect(mapState, mapDispatch)
type CurrentUserProps = ConnectedProps<typeof connector> & RouteComponentProps

export class CurrentUser extends Component<CurrentUserProps> {
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
      <Flex align="center">
        <Avatar size="sm" name={currentUser.first_name} />
        <Flex flexDir="column" marginLeft={4}>
          <Heading as="h3" size="sm">{currentUser.last_name}</Heading>
          <Text color="gray">Admin</Text>
        </Flex>
      </Flex>
    )
  }
}

export default connector(CurrentUser as any)
