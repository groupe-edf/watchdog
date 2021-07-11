import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Input,
  Stack,
  Text} from "@chakra-ui/react";
import React from "react";
import { Component } from "react";
import { ConnectedProps, connect } from "react-redux";
import { RouteComponentProps, withRouter } from "react-router-dom";
import { Card } from "../../components/Card";
import integrationsService from "../../services/integration"
import { ApplicationState } from "../../store";
import { IntegrationActionTypes } from "../../store/integrations/types";

export class CreateIntegration extends Component<any> {
  static INITIAL_STATE = {
    instance_url: '',
    instance_name: '',
    api_token: ''
  }
  constructor(props: any) {
    super(props);
    this.state = {
      ...CreateIntegration.INITIAL_STATE,
      error: ''
    };
  }
  static propKey(propertyName: string, value: any): object {
    return { [propertyName]: value };
  }
  handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    this.setStateWithEvent(event, event.target.name);
  }
  handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {

  }
  setStateWithEvent(event: any, columnType: string): void {
    this.setState(CreateIntegration.propKey(columnType, (event.target as any).value));
  }
  render() {
    return (
      <Box
        minH="100vh"
        paddingY="12"
        paddingX={{ base: '4', lg: '8' }}>
        <Box maxWidth="md" marginX="auto">
          <Card>
            <Stack spacing={6}>
            <form onSubmit={this.handleSubmit}>
              <FormControl isRequired>
                <FormLabel htmlFor="instance_url">Instance URL</FormLabel>
                <Input type="instance_url" name="instance_url" onChange={this.handleChange} />
              </FormControl>
              <FormControl isRequired>
                <FormLabel htmlFor="instance_name">Name your personal access token</FormLabel>
                <Input type="instance_name" name="instance_name" onChange={this.handleChange} />
              </FormControl>
              <FormControl isRequired>
                <FormLabel htmlFor="api_token">Personal access token (with api scope)</FormLabel>
                <Input type="api_token" name="api_token" onChange={this.handleChange} />
              </FormControl>
              <Button width="full" mt={4}
                type="submit"
                colorScheme="brand">
                Add
              </Button>
            </form>
            </Stack>
          </Card>
        </Box>
      </Box>
    )
  }
}

const mapState = (state: ApplicationState) => ({
  integration: state.integrations.integration,
})
const mapDispatch = {
  getIntegration: (payload: any) => ({ type: IntegrationActionTypes.INTEGRATION_FIND_BY_ID, payload }),
}
interface IntegrationParams {
  integrationId: string
}
const connector = connect(mapState, mapDispatch)
type IntegrationProps = ConnectedProps<typeof connector> & RouteComponentProps<IntegrationParams>

export class Integration extends Component<IntegrationProps> {
  constructor(props: IntegrationProps) {
    super(props);
  }
  componentDidMount() {
    const { match, getIntegration } = this.props
    integrationsService.findById(match.params.integrationId).then(response => {
      getIntegration(response.data)
    })
  }
  render() {
    return (
      <Box>
        <Text>Integration</Text>
      </Box>
    )
  }
}

export default withRouter(connector(Integration));
