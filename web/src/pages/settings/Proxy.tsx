import { Stack, FormControl, FormLabel, InputGroup, Input, Textarea, Button } from "@chakra-ui/react"

const Proxy = () => {
  return (
    <form>
      <Stack spacing={4} width="50%">
        <h1>HTTP Proxy Configuration</h1>
        <FormControl isRequired>
          <FormLabel htmlFor="host">Host</FormLabel>
          <InputGroup>
            <Input type="text" name="host" placeholder="Host" />
          </InputGroup>
        </FormControl>
        <FormControl isRequired>
          <FormLabel htmlFor="port">Port</FormLabel>
          <InputGroup>
            <Input type="text" name="port" placeholder="Port" />
          </InputGroup>
        </FormControl>
        <FormControl>
          <FormLabel htmlFor="username">Username</FormLabel>
          <InputGroup>
            <Input type="text" name="username" placeholder="Username" />
          </InputGroup>
        </FormControl>
        <FormControl>
          <FormLabel htmlFor="password">Password</FormLabel>
          <InputGroup>
            <Input type="password" name="password" placeholder="Password" />
          </InputGroup>
        </FormControl>
        <FormControl>
          <FormLabel htmlFor="no_proxy">No Proxy</FormLabel>
          <InputGroup>
            <Textarea name="no_proxy" placeholder="No Proxy" />
          </InputGroup>
        </FormControl>
      </Stack>
      <Stack marginTop={4} isInline={true}>
        <Button
          type="submit"
          loadingText="Validation.."
          colorScheme="brand"
          fontWeight="md">
          Validate Proxy
        </Button>
      </Stack>
    </form>
  )
}

export { Proxy }
