import {
  Box,
  Button,
  Stack,
  Switch,
  Table,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
  Text,
  useColorModeValue,
  Badge,
  Flex,
  Icon
} from "@chakra-ui/react";
import { Component, FC } from "react";
import { IoBuildOutline, IoMenu } from "react-icons/io5";
import { connect, ConnectedProps } from "react-redux";
import { Switch as SwitchRoute, Route, RouteComponentProps, withRouter } from "react-router-dom";
import { ButtonLink } from "../../components/ButtonLink";
import { withStatusIndicator } from "../../components/withStatusIndicator";
import policiesService from "../../services/policy"
import { ApplicationState } from "../../store";
import { Policy, PolicyActionTypes } from "../../store/policies/types";
import CreatePolicy from "./create";
import EditPolicy from "./edit";

interface PoliciesContentProps {
  data: Policy[];
}

export const PoliciesContent: FC<PoliciesContentProps> = ({ data }) => {
  const header = ['Name', 'Type', 'Conditions', 'State', ''];
  const togglePolicy = (id: number, enabled: boolean, event: React.ChangeEvent<HTMLInputElement>) => {
    policiesService.toggle(id, !enabled).then(data => {
    })
  };
  return (
    <Table variant="simple"
      background={useColorModeValue('white', 'gray.800')}>
      <Thead>
        <Tr>
          {header.map((value) => (
            <Th key={value}>{value}</Th>
          ))}
        </Tr>
      </Thead>
      <Tbody>
        {data && data.map(function(policy){
          return (
            <Tr>
              <Td>
                {policy.display_name}
                {policy.description && (
                  <Text fontSize={'sm'} color="gray">
                    {policy.description}
                  </Text>
                )}
              </Td>
              <Td><Badge>{policy.type}</Badge></Td>
              <Td>
                <Flex alignItems="center">
                  <Icon as={IoMenu} marginRight={2}/>
                  <Text fontWeight="bold">{policy.conditions?.length}</Text>
                </Flex>
              </Td>
              <Td align="right"><Switch defaultChecked={policy.enabled} colorScheme="brand" onChange={(event) => togglePolicy(policy.id, policy.enabled, event)}/></Td>
              <Td textAlign="right">
                <ButtonLink rightIcon={<IoBuildOutline />} to={`/policies/${policy.id}/edit`} colorScheme="gray" size="sm" variant="outline">
                  Edit
                </ButtonLink>
              </Td>
            </Tr>
          )
        })}
      </Tbody>
    </Table>
  )
}
PoliciesContent.displayName = 'Policies';
const PoliciesWithStatusIndicator = withStatusIndicator(PoliciesContent);

const mapState = (state: ApplicationState) => ({
  state: state.policies
})
const mapDispatch = {
  getPolicies: (payload: any) => ({ type: PolicyActionTypes.POLICIES_FIND_ALL, payload }),
}
const connector = connect(mapState, mapDispatch)
type PoliciesProps = ConnectedProps<typeof connector> & RouteComponentProps

class PoliciesList extends Component<PoliciesProps> {
  constructor(props: PoliciesProps) {
    super(props);
  }
  componentDidMount() {
    const { state, getPolicies } = this.props
    policiesService.findAll().then(response => {
      getPolicies(response.data)
    })
  }
  render() {
    const { match, state } = this.props
    return (
      <SwitchRoute>
        <Route path={`${match.url}/create`} component={CreatePolicy}/>
        <Route path={`${match.url}/:policyId/edit`} component={EditPolicy}/>
        <Route exact path={match.url}>
          <Box paddingBottom={4}>
            <Stack
              justify={'flex-end'}
              direction={'row'}>
              <ButtonLink to={`${match.url}/create`}
                colorScheme="brand"
                fontWeight="md">
                New
              </ButtonLink>
            </Stack>
          </Box>
          <PoliciesWithStatusIndicator data={state.policies}/>
        </Route>
      </SwitchRoute>
    )
  }
}

export default withRouter(connector(PoliciesList));
