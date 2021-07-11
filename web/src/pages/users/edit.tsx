import { Stack, FormControl, FormLabel, Input, Box, Button } from "@chakra-ui/react"
import { Component } from "react"
import { connect, ConnectedProps } from "react-redux"
import { RouteComponentProps, withRouter } from "react-router-dom"
import { ApplicationState } from "../../store"
import { UserActionTypes, User } from "../../store/users/types"

const mapState = (state: ApplicationState) => ({
  currentUser: state.authentication.currentUser,
  state: state.users
})
const mapDispatch = {
  getUsers: (payload: any) => ({ type: UserActionTypes.USERS_FIND_ALL, payload }),
}
const connector = connect(mapState, mapDispatch)
type UserEditProps = ConnectedProps<typeof connector> & RouteComponentProps & {
  user: User
}
export class UserEdit extends Component<UserEditProps> {
  render() {
    return (
      <Box padding={6} paddingX={{ base: '4', md: '6' }} background="white">
        <Stack paddingBottom={5} spacing={4}>
          <FormControl>
            <FormLabel>Display Name</FormLabel>
            <Input type="text" name="display_name" value="" />
          </FormControl>
        </Stack>
        <Box paddingX={{ base: 4, sm: 6 }} paddingY={3} background="gray.50" textAlign="right">
          <Button
            type="submit"
            colorScheme="brand"
            loadingText="Updating.."
            _focus={{ shadow: "" }}
            fontWeight="md">
            Update
          </Button>
        </Box>
      </Box>
    )
  }
}

export default withRouter(connector(UserEdit));
