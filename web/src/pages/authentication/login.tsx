import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Heading,
  Input,
  Stack,
  Text,
  Link,
  SimpleGrid,
  VisuallyHidden,
  Alert,
  AlertDescription,
  AlertIcon,
  Checkbox
} from "@chakra-ui/react";
import { Component } from "react";
import { IoLogoGitlab } from "react-icons/io5";
import { connect, ConnectedProps } from "react-redux";
import { Link as RouterLink, RouteComponentProps } from "react-router-dom";
import { Card } from "../../components/Card";
import { DividerWithText } from "../../components/DividerWithText";
import authenticationService from "../../services/authentication";
import { ApplicationState } from "../../store";
import { AuthenticationActionTypes } from "../../store/authentication/types";

const mapState = (state: ApplicationState) => ({
  currentUser: state.authentication.currentUser,
  settings: state.global.settings
})
const mapDispatch = {
  login: (payload: any) => ({ type: AuthenticationActionTypes.AUTHENTICATION_LOGIN, payload }),
}
const connector = connect(mapState, mapDispatch)
type LoginProps = ConnectedProps<typeof connector> & RouteComponentProps

interface LoginState {
  email: string
  error: string
  isLoading: boolean
  password: string
}

export class Login extends Component<LoginProps, LoginState> {
  static INITIAL_STATE = {
    email: '',
    password: ''
  };
  constructor(props: LoginProps) {
    super(props);
    this.state = {
      ...Login.INITIAL_STATE,
      error: '',
      isLoading: false
    };
  }
  static propKey(propertyName: string, value: any): object {
    return { [propertyName]: value };
  }
  setStateWithEvent(event: any, columnType: string): void {
    this.setState(Login.propKey(columnType, (event.target as any).value));
  }
  onSubmit = (event: any) => {
    event.preventDefault()
    this.setState({ isLoading: true })
    const { login, history } = this.props
    const { email, password } = this.state;
    authenticationService.login({email: email, password: password}).then(response => {
      localStorage.setItem('user', JSON.stringify({
        email: response.data.email,
        first_name: response.data.first_name,
        last_name: response.data.last_name
      }))
      if (response.data.token) {
        localStorage.setItem('token', response.data.token)
        history.push('/')
        login(response.data)
      }
    }).catch(response => {
      this.setState({ error: response.error.detail})
    }).finally(() => {
      this.setState({ isLoading: false })
    })
  }
  render() {
    const { settings } = this.props
    const { isLoading, email, password } = this.state
    return (
      <Box
        background="gray.50"
        minH="100vh"
        paddingY="12"
        paddingX={{ base: '4', lg: '8' }}>
        <Box maxWidth="md" marginX="auto">
          <Heading textAlign="center" size="xl" fontWeight="bold" marginBottom="4">
            Sign in to your account
          </Heading>
          {settings.enable_signup && <Text align="center" maxW="md" fontWeight="medium">
            <Text as="span">Don't have an account?</Text>
            <Link as={RouterLink} to="/register" color="brand.100" marginStart="1" display={{ base: 'block', sm: 'inline' }}>Register</Link>
          </Text>
          }
          <Card marginTop="8">
            <Stack spacing={6}>
              <form onSubmit={event => this.onSubmit(event)}>
                {this.state.error &&
                <Alert status="error" marginBottom={4}>
                  <AlertIcon />
                  <AlertDescription>{this.state.error}</AlertDescription>
                </Alert>
                }
                <FormControl isRequired>
                  <FormLabel htmlFor='email'>Email</FormLabel>
                  <Input type="text" id="email" name="email" value={email} onChange={event => this.setStateWithEvent(event, "email")} />
                </FormControl>
                <FormControl mt={6}>
                  <FormLabel>Password</FormLabel>
                  <Input type="password" name="password" value={password} onChange={event => this.setStateWithEvent(event, "password")} />
                </FormControl>
                <Stack spacing={6} marginTop={2}>
                  <Stack
                    direction={{ base: 'column', sm: 'row' }}
                    align={'start'}
                    justify={'space-between'}>
                    <Checkbox isChecked={true} colorScheme="brand">Remember me</Checkbox>
                    <Link as={RouterLink} to="/reset" color={'brand.100'}>Forgot password ?</Link>
                  </Stack>
                  <Button width="full" type="submit" colorScheme="brand" isLoading={isLoading}>
                    Login
                  </Button>
                </Stack>
              </form>
            </Stack>
            {settings.enable_oauth_signup &&
            <>
              <DividerWithText mt="6">or continue with</DividerWithText>
              <SimpleGrid marginTop="6" columns={1} spacing="1">
                <Button color="currentColor" variant="outline">
                  <VisuallyHidden>Login with Gitlab</VisuallyHidden>
                  <IoLogoGitlab />
                </Button>
              </SimpleGrid>
            </>
            }
          </Card>
        </Box>
      </Box>
    )
  }
}

export default connector(Login);
