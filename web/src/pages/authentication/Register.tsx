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
  Text,
  useColorModeValue
} from "@chakra-ui/react"
import { Card, CardBody } from "@saas-ui/react"
import { useState } from "react"
import { useDispatch, useSelector } from "react-redux"
import { Link as RouterLink, useNavigate } from "react-router-dom"
import { AppDispatch, RootState } from "../../configureStore"
import { register } from "../../store/slices/authentication"

interface RegisterState {
  email: string
  error: string
  first_name: string
  isLoading: boolean
  last_name: string
  password: string
}
const Register = (props: any) => {
  const dispatch = useDispatch<AppDispatch>()
  const navigate = useNavigate()
  const { message, settings } = useSelector((state: RootState) => state.global)
  const [values, setValues] = useState({})
  const [loading, setLoading] = useState(false)
  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const name = event.target.name;
    const value = event.target.value;
    setValues(values => ({...values, [name]: value}))
  }
  const onSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setLoading(true)
    dispatch(register(values))
      .unwrap()
      .then(() => {
        navigate("/login")
      })
      .catch(() => {
        setLoading(false)
      })
  }
  return (
    <Box
      background={useColorModeValue('gray.100', 'gray.900')}
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
          <CardBody>
            <Stack spacing={6}>
              <form onSubmit={onSubmit}>
                {message &&
                <Alert status="error">
                  <AlertIcon />
                  <AlertDescription>{message.detail}</AlertDescription>
                </Alert>
                }
                <FormControl isRequired>
                  <FormLabel htmlFor="email">Email</FormLabel>
                  <Input type="email" name="email" onChange={handleChange} />
                </FormControl>
                <FormControl isRequired>
                  <FormLabel htmlFor="first_name">Firstname</FormLabel>
                  <Input type="text" name="first_name" onChange={handleChange} />
                </FormControl>
                <FormControl isRequired>
                  <FormLabel htmlFor="last_name">Lastname</FormLabel>
                  <Input type="text" name="last_name" onChange={handleChange} />
                </FormControl>
                <FormControl isRequired>
                  <FormLabel htmlFor="password">Password</FormLabel>
                  <Input type="password" name="password" onChange={handleChange} />
                </FormControl>
                <Button
                  width="full"
                  marginTop={4}
                  type="submit"
                  isLoading={loading}
                  colorScheme="brand">
                  Register
                </Button>
              </form>
            </Stack>
          </CardBody>
        </Card>
      </Box>
    </Box>
  )
}
export default Register
