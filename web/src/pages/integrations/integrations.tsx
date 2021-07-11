import {
  Box,
  Button,
  ButtonGroup,
  HStack,
  Icon,
  IconButton,
  Link,
  Stack,
  Table,
  Tbody,
  Td,
  Text,
  Th,
  Thead,
  Tooltip,
  Tr
} from "@chakra-ui/react";
import React, { Component } from "react";
import { ConnectedProps, connect } from "react-redux";
import { Link as ReactRouterLink, RouteComponentProps, withRouter } from "react-router-dom";
import { ApplicationState } from "../../store";
import { IntegrationActionTypes } from "../../store/integrations/types";
import integrationsService from "../../services/integration"
import { IoWarningOutline, IoSyncOutline, IoTrashOutline, IoBatteryDeadOutline, IoFlashOffOutline } from "react-icons/io5";
import { AddIntegration } from "../../components/integrations/AddIntegration";

const mapState = (state: ApplicationState) => ({
  state: state.integrations,
})
const mapDispatch = {
  getIntegrations: (payload: any) => ({ type: IntegrationActionTypes.INTEGRATION_FIND_ALL, payload }),
  synchronize: (payload: any) => ({ type: IntegrationActionTypes.INTEGRATION_SYNCHROONIZE, payload }),
}
const connector = connect(mapState, mapDispatch)
type IntegrationProps = ConnectedProps<typeof connector> & RouteComponentProps

export class IntegrationsList extends Component<IntegrationProps, {
  isSubmitting: boolean
}> {
  constructor(props: IntegrationProps) {
    super(props);
    this.state = {
      isSubmitting: false
    }
  }
  componentDidMount() {
    const { getIntegrations } = this.props
    integrationsService.findAll().then(response => {
      getIntegrations(response.data)
    })
  }
  render() {
    const header = ['Instance Name', 'Instance URL', 'Created At', 'Last Sync', 'Actions']
    const { match, state } = this.props
    const { isSubmitting } = this.state
    const synchronize = async (integrationId: string, event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
      const { synchronize } = this.props
      event.preventDefault();
      this.setState({isSubmitting: true})
      integrationsService.synchronize(integrationId).then(response => {
        synchronize(response.data)
        this.setState({isSubmitting: false})
      })
    }
    return (
      <>
      <Box paddingBottom={4}>
        <Stack
          justify={'flex-end'}
          direction={'row'}>
          <AddIntegration/>
        </Stack>
      </Box>
      <Table variant="simple" background="white">
        <Thead>
          <Tr>
            {header.map((value) => (
              <Th key={value}>{value}</Th>
            ))}
          </Tr>
        </Thead>
        <Tbody>
          {state.integrations.length > 0 ? (
            state.integrations.map(function(integration){
            return (
              <Tr key={integration.id}>
                <Td>
                  <Link as={ReactRouterLink} color="brand.100" to={`${match.url}/integrations/${integration.id}`} style={{ textDecoration: 'none' }}>
                    {integration.instance_name}
                  </Link>
                </Td>
                <Td>
                  {integration.instance_url}
                </Td>
                <Td>
                  {integration.created_at && new Intl.DateTimeFormat("en-GB", {
                    year: "numeric",
                    month: "long",
                    day: "2-digit",
                    hour: "2-digit",
                    minute: "2-digit",
                    second: "2-digit",
                  }).format(Date.parse(integration.created_at))}
                </Td>
                <Td>
                  <HStack>
                    <Text>
                      {integration.synced_at && new Intl.DateTimeFormat("en-GB", {
                        year: "numeric",
                        month: "long",
                        day: "2-digit",
                        hour: "2-digit",
                        minute: "2-digit",
                        second: "2-digit",
                      }).format(Date.parse(integration.synced_at))}
                    </Text>
                    {integration.syncing_error &&
                    <Tooltip label={integration.syncing_error}>
                      <Icon as={IoWarningOutline} color="brand.100" />
                    </Tooltip>
                    }
                  </HStack>
                </Td>
                <Td>
                <ButtonGroup size="sm" isAttached colorScheme="brand">
                  <Button
                    leftIcon={<IoSyncOutline />}
                    isLoading={isSubmitting}
                    loadingText="Synchronizing"
                    onClick={(event) => synchronize(integration.id, event)}>
                    Synchronize
                  </Button>
                  <IconButton aria-label="Delete" icon={<IoTrashOutline />}/>
                </ButtonGroup>
                </Td>
              </Tr>
            )
          })) : (
            <Tr>
              <Td colSpan={header.length} textAlign="center" color="grey" paddingX={4}>
                  <Icon fontSize="64" as={IoFlashOffOutline} />
                  <Text marginTop={4}>No integrations found</Text>
                </Td>
            </Tr>
          )}
        </Tbody>
      </Table>
      </>
    )
  }
}

export default withRouter(connector(IntegrationsList));
