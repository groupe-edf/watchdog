import {
  Button,
  FormControl,
  SimpleGrid,
  FormLabel,
  Input,
  InputGroup,
  Stack,
  Tab,
  Table,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
  Text,
  Box,
  GridItem,
  Heading
} from "@chakra-ui/react"
import { Component } from "react"
import { connect, ConnectedProps } from "react-redux"
import { RouteComponentProps } from "react-router-dom"
import { API_PATH } from "../constants"
import { fetchData } from "../services/commons"
import { ApplicationState } from "../store"
import { AuthenticationActionTypes } from "../store/authentication/types"
import userService from "../services/user"
import { APIKeyList } from "./api"

const mapState = (state: ApplicationState) => ({
  currentUser: state.authentication.currentUser
})
const mapDispatch = {
  getCurrentUser: (payload: any) => ({ type: AuthenticationActionTypes.AUTHENTICATION_CURRENT_USER, payload }),
}
const connector = connect(mapState, mapDispatch)
type ProfileProps = ConnectedProps<typeof connector> & RouteComponentProps

export class Profile extends Component<ProfileProps, any> {
  static navigationOptions = {
    title: 'Profile'
  }
  constructor(props: ProfileProps) {
    super(props)
    this.state = {
      current_password: "",
      password: "",
      confirm_password: ""
    }
  }
  componentDidMount() {
    const { currentUser, getCurrentUser } = this.props
    if (currentUser.email === "") {
      fetchData("GET", `${API_PATH}/profile`).then(response => {
        getCurrentUser(response.data)
      })
    }
  }
  static propKey(propertyName: string, value: any): object {
    return { [propertyName]: value };
  }
  setStateWithEvent(event: any, columnType: string): void {
    this.setState(Profile.propKey(columnType, (event.target as any).value));
  }
  handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    this.setStateWithEvent(event, event.target.name);
  }
  changePassword = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    userService.changePassword({
      current_password: this.state.current_password,
      password: this.state.password,
      confirm_password: this.state.confirm_password,
    }).then(response => {})
  }
  updateProfile = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
  }
  render() {
    const { currentUser } = this.props
    return (
      <Tabs>
        <TabList>
          <Tab>Profile</Tab>
          <Tab isDisabled={currentUser.provider === "ldap"}>Security</Tab>
          <Tab>API</Tab>
        </TabList>
        <TabPanels background="white" borderBottomRadius="md">
          <TabPanel>
            <SimpleGrid
              display={{ base: "initial", md: "grid" }}
              columns={{ md: 3 }}
              spacing={{ md: 6 }}>
              <GridItem colSpan={{ md: 1 }}>
                <Box px={[4, 0]}>
                  <Heading fontSize="lg" fontWeight="medium" lineHeight="6">
                    Personal Information
                  </Heading>
                  <Text
                    marginTop={1}
                    fontSize="sm"
                    color="gray.600">
                    Use a permanent address where you can receive mail.
                  </Text>
                </Box>
              </GridItem>
              <GridItem mt={[5, null, 0]} colSpan={{ md: 2 }}>
              <form onSubmit={this.updateProfile}>
                <Stack
                  paddingY={5}
                  background="white"
                  spacing={6}>
                  <SimpleGrid columns={6} spacing={4}>
                    <FormControl as={GridItem} isRequired colSpan={[6, 3]}>
                      <FormLabel htmlFor="first_name">First Name</FormLabel>
                      <InputGroup>
                        <Input
                          type="first_name"
                          name="first_name"
                          value={currentUser.first_name}
                          onChange={this.handleChange}
                          focusBorderColor="brand.100"
                          placeholder="First Name" />
                      </InputGroup>
                    </FormControl>
                    <FormControl as={GridItem} isRequired colSpan={[6, 3]}>
                      <FormLabel htmlFor="last_name">Last Name</FormLabel>
                      <InputGroup>
                        <Input
                          type="last_name"
                          name="last_name"
                          value={currentUser.last_name}
                          onChange={this.handleChange}
                          focusBorderColor="brand.100"
                          placeholder="Last Name" />
                      </InputGroup>
                    </FormControl>
                    <FormControl as={GridItem} isRequired colSpan={[6, 6]}>
                      <FormLabel htmlFor="email">Email</FormLabel>
                      <InputGroup>
                        <Input
                          type="email"
                          name="email"
                          value={currentUser.email}
                          onChange={this.handleChange}
                          placeholder="Email"
                          isDisabled={true} />
                      </InputGroup>
                    </FormControl>
                  </SimpleGrid>
                </Stack>
                <Box
                  paddingX={{ base: 4, sm: 6 }}
                  paddingY={3}
                  background="gray.50"
                  textAlign="right">
                  <Button
                    type="submit"
                    colorScheme="brand"
                    loadingText="Fetching.."
                    _focus={{ shadow: "" }}
                    fontWeight="md">
                    Update
                  </Button>
                </Box>
              </form>
              </GridItem>
            </SimpleGrid>
          </TabPanel>
          <TabPanel>
            <form onSubmit={this.changePassword}>
              <Stack spacing={4} width="50%">
                <FormControl isRequired>
                  <FormLabel htmlFor="current_password">Current Password</FormLabel>
                  <InputGroup>
                    <Input type="password" name="current_password" placeholder="Current password" onChange={this.handleChange} />
                  </InputGroup>
                </FormControl>
                <FormControl isRequired>
                  <FormLabel htmlFor="password">Password</FormLabel>
                  <InputGroup>
                    <Input type="password" name="password" placeholder="Password" onChange={this.handleChange} />
                  </InputGroup>
                </FormControl>
                <FormControl isRequired>
                  <FormLabel htmlFor="confirm_password">Confirm Password</FormLabel>
                  <InputGroup>
                    <Input type="password" name="confirm_password" placeholder="Confirm password" onChange={this.handleChange} />
                  </InputGroup>
                </FormControl>
              </Stack>
              <Stack marginTop={4} isInline={true}>
                <Button
                  type="submit"
                  loadingText="Fetching.."
                  colorScheme="brand"
                  fontWeight="md">
                  Update
                </Button>
              </Stack>
            </form>
          </TabPanel>
        </TabPanels>
      </Tabs>
    )
  }
}

export default connector(Profile);
