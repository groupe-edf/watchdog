import {
  Box,
  Button,
  Flex,
  FormControl,
  FormLabel,
  Input,
  Stack,
  useColorModeValue
} from "@chakra-ui/react";
import { useState } from "react";
import { API_PATH } from "../constants";

export default function Login({ setToken }: { setToken: CallableFunction }) {
  const initialState = {
    username: "",
    password: ""
  };
  const [isSubmitting, setSubmitting] = useState(false);
  const [values, setValues] = useState(initialState);
  const onChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setValues({ ...values, [event.target.name]: event.target.value });
  };
  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const token = await login(values);
  }
  async function login(credentials: {username: string|undefined, password: string|undefined}) {
    fetch(`${API_PATH}/login`, {
      method: 'POST',
      body: JSON.stringify(credentials),
      headers: { 'Content-Type': 'application/json' }
    })
    .then(response => response.json())
  }
  return (
    <Flex width="full" align="center" justifyContent="center">
      <Box p={2}>
        <Box my={4} textAlign="left">
          <form onSubmit={handleSubmit}>
            <FormControl>
              <FormLabel>Unsername</FormLabel>
              <Input type="text" name="username" onChange={onChange} />
            </FormControl>
            <FormControl mt={6}>
              <FormLabel>Password</FormLabel>
              <Input type="password" name="password" onChange={onChange} />
            </FormControl>
            <Button width="full" mt={4}
              type="submit"
              isLoading={isSubmitting}
              colorScheme="teal">
              Sign In
            </Button>
          </form>
        </Box>
      </Box>
    </Flex>
  )
}

export { Login }
