import {
  Box,
  Checkbox,
  Flex,
  FormControl,
  FormHelperText,
  FormLabel,
  Stack,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
  Text,
  useToast
} from "@chakra-ui/react"
import { Component, createContext } from "react"
import { connect, ConnectedProps } from "react-redux"
import { Route, RouteComponentProps, Switch as SwitchRoute, withRouter } from "react-router-dom"
import integrationsService from "../../services/integration"
import { ApplicationState } from "../../store"
import { IntegrationActionTypes } from "../../store/integrations/types"
import APIKeyList from "../api"
import { Proxy } from "./proxy"

const ToastContext = createContext(() => {});
function ToastProvider({ children }: any) {
  const toast = useToast();
  return (
    <ToastContext.Provider value={toast}>{children}</ToastContext.Provider>
  );
}

const mapState = (state: ApplicationState) => ({
  integrations: state.integrations,
  settings: state.global.settings
})
const mapDispatch = {
  getIntegrations: (payload: any) => ({ type: IntegrationActionTypes.INTEGRATION_FIND_ALL, payload }),
}
const connector = connect(mapState, mapDispatch)
type SettingsProps = ConnectedProps<typeof connector> & RouteComponentProps

export class Settings extends Component<SettingsProps> {
  constructor(props: SettingsProps) {
    super(props);
  }
  componentDidMount() {
    const { integrations, getIntegrations } = this.props
    if (integrations.integrations.length === 0) {
      integrationsService.findAll().then(response => {
        getIntegrations(response.data)
      })
    }
  }
  render() {
    const { settings, history, location, match } = this.props
    return (
      <SwitchRoute>
        <Route exact path={match.url}>
          <Tabs>
            <TabList>
              <Tab>Global</Tab>
              <Tab>Proxy</Tab>
              <Tab>API Keys</Tab>
            </TabList>
            <TabPanels background="white" borderBottomRadius="md">
              <TabPanel>
                <Box
                  as="legend"
                  fontSize="md"
                  color="gray.900">
                  Security
                </Box>
                <Stack mt={4} spacing={4}>
                  <Flex alignItems="start">
                    <Flex alignItems="center" height={6}>
                      <Checkbox colorScheme="brand" isChecked={settings.enable_signup}/>
                    </Flex>
                    <FormControl marginLeft={2}>
                      <FormLabel marginBottom={0} fontSize="sm">Enable Signup</FormLabel>
                      <FormHelperText marginTop={0}>
                        Enable users to create a Watchdog account
                      </FormHelperText>
                    </FormControl>
                  </Flex>
                  <Flex alignItems="start">
                    <Flex alignItems="center" height={6}>
                      <Checkbox colorScheme="brand" isChecked={settings.enable_oauth_signup}/>
                    </Flex>
                    <FormControl marginLeft={2}>
                      <FormLabel marginBottom={0} fontSize="sm">Enable OAuth Signup</FormLabel>
                      <FormHelperText marginTop={0}>
                        Enable user to connect to Watchdog with several OAuth application such as Gitlab
                      </FormHelperText>
                    </FormControl>
                  </Flex>
                </Stack>
              </TabPanel>
              <TabPanel>
                <Proxy></Proxy>
              </TabPanel>
              <TabPanel>
                <APIKeyList></APIKeyList>
              </TabPanel>
            </TabPanels>
          </Tabs>
        </Route>
      </SwitchRoute>
    )
  }
}
export default withRouter(connector(Settings));
