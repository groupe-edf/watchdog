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
  Checkbox,
  InputGroup,
  InputRightElement,
  useColorModeValue
} from "@chakra-ui/react"
import { Card, CardBody } from "@saas-ui/react"
import { useEffect, useState } from "react"
import { useForm } from "react-hook-form"
import { IoLogoGitlab } from "react-icons/io5"
import { useDispatch, useSelector } from "react-redux"
import { useNavigate, Link as RouterLink } from "react-router-dom"
import { DividerWithText } from "../../components/DividerWithText"
import { AppDispatch, RootState } from "../../configureStore"
import { login } from "../../store/slices/authentication"
import { clearMessage } from "../../store/slices/global"

export type Credentials = {
  email: string
  password: string
}
const Login = (props: any) => {
  const navigate = useNavigate()
  const dispatch = useDispatch<AppDispatch>()
  const { message, settings } = useSelector((state: RootState) => state.global)
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<Credentials>()
  useEffect(() => {
    dispatch(clearMessage());
  }, [dispatch])
  const [show, setShow] = useState(false)
  const handleTogglePassword = () => setShow(!show)
  const onSubmit = (values: Credentials) => {
    dispatch(login(values))
      .unwrap()
      .then(() => {
        navigate("/")
      })
  }
  return (
    <Box
      background={useColorModeValue('gray.100', 'gray.900')}
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
          <CardBody>
          <Stack spacing={6}>
            <form onSubmit={handleSubmit(onSubmit)}>
              {message && (
                <Alert status="error" marginBottom={4}>
                  <AlertIcon />
                  <AlertDescription>{message.detail}</AlertDescription>
                </Alert>
              )}
              <FormControl isRequired>
                <FormLabel htmlFor='email'>Email</FormLabel>
                <Input type="text" {...register('email', {
                  required: 'Email is required'
                })} />
              </FormControl>
              <FormControl mt={6}>
                <FormLabel>Password</FormLabel>
                <InputGroup>
                  <Input type={show ? "text" : "password"} {...register('password', {
                    required: 'Password is required'
                  })}/>
                  <InputRightElement width="4.5rem">
                    <Button h="1.75rem" size="sm" onClick={handleTogglePassword}>
                      {show ? "Hide" : "Show"}
                    </Button>
                  </InputRightElement>
                </InputGroup>
              </FormControl>
              <Stack spacing={6} marginTop={2}>
                <Stack
                  direction={{ base: 'column', sm: 'row' }}
                  align={'start'}
                  justify={'space-between'}>
                  <Checkbox isChecked={true} colorScheme="brand">Remember me</Checkbox>
                  <Link as={RouterLink} to="/reset" color={'brand.100'}>Forgot password ?</Link>
                </Stack>
                <Button width="full" type="submit" colorScheme="brand" isLoading={isSubmitting}>
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
          </CardBody>
        </Card>
      </Box>
    </Box>
  )
}
export default Login
