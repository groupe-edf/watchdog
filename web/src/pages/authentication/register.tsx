import {
  Alert,
  AlertDescription,
  AlertIcon,
  Box,
  Button,
  FormControl,
  FormLabel,
  Heading,
  Input,
  Link,
  Stack,
  Text
} from "@chakra-ui/react";
import { Component } from "react";
import { Link as RouterLink, RouteComponentProps } from "react-router-dom";
import { Card } from "../../components/Card";
import authenticationService from "../../services/authentication"

interface RegisterState {
  email: string
  error: string
  first_name: string
  isLoading: boolean
  last_name: string
  password: string
}

export class Register extends Component<RouteComponentProps, RegisterState> {
  static INITIAL_STATE = {
    email: '',
    first_name: '',
    last_name: '',
    password: ''
  }
  constructor(props: any) {
    super(props);
    this.state = {
      ...Register.INITIAL_STATE,
      error: '',
      isLoading: false
    };
  }
  static propKey(propertyName: string, value: any): object {
    return { [propertyName]: value };
  }
  setStateWithEvent(event: any, columnType: string): void {
    this.setState(Register.propKey(columnType, (event.target as any).value));
  }
  handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    this.setStateWithEvent(event, event.target.name);
  }
  handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    this.setState({ isLoading: true })
    const { history } = this.props
    authenticationService.register(this.state).then(response => {
      history.push('/login');
    }).catch(response => {
      this.setState({ error: response.error.detail})
    }).finally(() => {
      this.setState({ isLoading: false })
    })
  }
  render() {
    const { isLoading } = this.state
    return (
      <Box
        background="gray.50"
        minH="100vh"
        paddingY="12"
        paddingX={{ base: '4', lg: '8' }}>
        <Box maxWidth="md" marginX="auto">
          <Heading textAlign="center" size="xl" fontWeight="bold">
            Create your account
          </Heading>
          <Text mt="4" mb="8" align="center" maxW="md" fontWeight="medium">
            <Text as="span">Already have an account?</Text>
            <Link as={RouterLink} to="/login" color="brand.100" marginStart="1" display={{ base: 'block', sm: 'inline' }}>Login</Link>
          </Text>
          <Card>
            <Stack spacing={6}>
              <form onSubmit={this.handleSubmit}>
                {this.state.error &&
                <Alert status="error">
                  <AlertIcon />
                  <AlertDescription>{this.state.error}</AlertDescription>
                </Alert>
                }
                <FormControl isRequired>
                  <FormLabel htmlFor="email">Email</FormLabel>
                  <Input type="email" name="email" onChange={this.handleChange} />
                </FormControl>
                <FormControl isRequired>
                  <FormLabel htmlFor="first_name">Firstname</FormLabel>
                  <Input type="text" name="first_name" onChange={this.handleChange} />
                </FormControl>
                <FormControl isRequired>
                  <FormLabel htmlFor="last_name">Lastname</FormLabel>
                  <Input type="text" name="last_name" onChange={this.handleChange} />
                </FormControl>
                <FormControl isRequired>
                  <FormLabel htmlFor="password">Password</FormLabel>
                  <Input type="password" name="password" onChange={this.handleChange} />
                </FormControl>
                <Button
                  width="full"
                  marginTop={4}
                  type="submit"
                  isLoading={isLoading}
                  colorScheme="brand">
                  Register
                </Button>
              </form>
            </Stack>
          </Card>
        </Box>
      </Box>
    )
  }
}
