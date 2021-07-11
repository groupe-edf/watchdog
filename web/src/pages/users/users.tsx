import { Component } from "react";
import { connect, ConnectedProps } from "react-redux";
import { Switch as SwitchRoute, Link as RouterLink, Route, RouteComponentProps, withRouter } from "react-router-dom";
import { Table, Tbody, Td, Th, Thead, Tr, Switch, Badge, Flex, Icon, Text, Menu, MenuButton, MenuItem, MenuList, Input, InputGroup, InputLeftElement } from "@chakra-ui/react";
import { IoEllipsisVertical, IoLinkOutline, IoMailOutline, IoSearchOutline } from "react-icons/io5";
import { ApplicationState } from "../../store";
import { User, UserActionTypes } from "../../store/users/types";
import userService from "../../services/user"
import { UserEdit } from "./edit";

const mapState = (state: ApplicationState) => ({
  currentUser: state.authentication.currentUser,
  state: state.users
})
const mapDispatch = {
  getUsers: (payload: any) => ({ type: UserActionTypes.USERS_FIND_ALL, payload }),
}
const connector = connect(mapState, mapDispatch)
type UserListProps = ConnectedProps<typeof connector> & RouteComponentProps & {
  user: User
}
export class UserList extends Component<UserListProps> {
  constructor(props: UserListProps) {
    super(props);
  }
  componentDidMount() {
    const { state, getUsers } = this.props
    if (state.users.length === 0) {
      userService.findAll().then(response => {
        getUsers(response.data)
      })
    }
  }
  render() {
    const header = ['Locked', 'Email', 'Full Name', 'Provider', 'Created At', 'Last Login', 'Role', ''];
    const { currentUser, match, state } = this.props
    return (
      <SwitchRoute>
        <Route path={`${match.url}/:userId/edit`} component={UserEdit}/>
        <Route exact path={match.url}>
        <Flex
          as="header"
          align="center"
          justify="space-between"
          marginBottom={4}
          width="full">
          <InputGroup width="96" display={{ base: "none", md: "flex" }}>
            <InputLeftElement children={<IoSearchOutline/>} />
            <Input name="query" placeholder="Search for users..."  background="white"/>
          </InputGroup>
        </Flex>
        <Table variant="simple" background="white">
          <Thead>
            <Tr>
              {header.map((value) => (
                <Th key={value}>{value}</Th>
              ))}
            </Tr>
          </Thead>
          <Tbody>
          {state.users && state.users.map(function(user){
            return (
              <Tr>
                <Td><Switch defaultChecked={user.locked} isReadOnly={true} colorScheme="brand"/></Td>
                <Td>
                  <Flex alignItems="center">
                    <Icon as={IoMailOutline} marginRight={2}/>
                    <Text>{user.email}</Text>
                    {user.email === currentUser.email &&
                      <Badge variant="outline" colorScheme="brand" size="xs" marginLeft={2}>Me</Badge>
                    }
                  </Flex>
                  {user.username &&
                    <Flex alignItems="center">
                      <Icon as={IoLinkOutline} marginRight={2}/>
                      <Text>{user.username}</Text>
                    </Flex>
                  }
                </Td>
                <Td>{user.first_name} {user.last_name}</Td>
                <Td><Badge variant="outline" colorScheme="brand">{user.provider}</Badge></Td>
                <Td>
                  {user.created_at && new Intl.DateTimeFormat("en-GB", {
                    year: "numeric",
                    month: "long",
                    day: "2-digit",
                    hour: "2-digit",
                    minute: "2-digit",
                    second: "2-digit",
                  }).format(Date.parse(user.created_at))}
                </Td>
                <Td>
                  {user.last_login && new Intl.DateTimeFormat("en-GB", {
                    year: "numeric",
                    month: "long",
                    day: "2-digit",
                    hour: "2-digit",
                    minute: "2-digit",
                    second: "2-digit",
                  }).format(Date.parse(user.last_login))}
                </Td>
                <Td>{user.role}</Td>
                <Td>
                  <Menu>
                    <MenuButton>
                    <IoEllipsisVertical />
                    </MenuButton>
                    <MenuList>
                      <MenuItem as={RouterLink} to={`/users/${user.id}/edit`}>Edit</MenuItem>
                      <MenuItem>Lock</MenuItem>
                    </MenuList>
                  </Menu>
                </Td>
              </Tr>
            )
          })}
          </Tbody>
        </Table>
        </Route>
      </SwitchRoute>
    )
  }
}

export default withRouter(connector(UserList));
