import { Button, FormControl, FormLabel, Input, InputGroup, Stack, Textarea } from "@chakra-ui/react";
import { Component } from "react";

export class Proxy extends Component<any> {
  handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    this.setStateWithEvent(event, event.target.name);
  }
  setStateWithEvent(event: any, columnType: string): void {
  }
  updateProxy = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
  }
  render() {
    return (
      <form onSubmit={this.updateProxy}>
        <Stack spacing={4} width="50%">
          <h1>HTTP Proxy Configuration</h1>
          <FormControl isRequired>
            <FormLabel htmlFor="host">Host</FormLabel>
            <InputGroup>
              <Input type="text" name="host" placeholder="Host" onChange={this.handleChange} />
            </InputGroup>
          </FormControl>
          <FormControl isRequired>
            <FormLabel htmlFor="port">Port</FormLabel>
            <InputGroup>
              <Input type="text" name="port" placeholder="Port" onChange={this.handleChange} />
            </InputGroup>
          </FormControl>
          <FormControl>
            <FormLabel htmlFor="username">Username</FormLabel>
            <InputGroup>
              <Input type="text" name="username" placeholder="Username" onChange={this.handleChange} />
            </InputGroup>
          </FormControl>
          <FormControl>
            <FormLabel htmlFor="password">Password</FormLabel>
            <InputGroup>
              <Input type="password" name="password" placeholder="Password" onChange={this.handleChange} />
            </InputGroup>
          </FormControl>
          <FormControl>
            <FormLabel htmlFor="no_proxy">No Proxy</FormLabel>
            <InputGroup>
              <Textarea type="text" name="no_proxy" placeholder="No Proxy" />
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
}
