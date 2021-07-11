import {
  Box,
  FormControl,
  FormLabel,
  GridItem,
  Heading,
  Input,
  Select,
  SimpleGrid,
  Stack,
  Switch,
  Table,
  Textarea,
  Th,
  Thead,
  Tr,
  Tbody,
  Td,
  IconButton,
  Button,
  Editable,
  EditableInput,
  EditablePreview
} from "@chakra-ui/react";
import { Component } from "react";
import { connect, ConnectedProps } from "react-redux";
import { RouteComponentProps, withRouter } from "react-router-dom";
import { ApplicationState } from "../../store";
import { PolicyActionTypes } from "../../store/policies/types";
import policiesService from "../../services/policy"
import { Card } from "../../components/Card";
import { IoTrashOutline } from "react-icons/io5";
import analysis from "../../services/analysis";
import { Pattern } from "../../components/Pattern";

const mapState = (state: ApplicationState) => ({
  policy: state.policies.policy,
})
const mapDispatch = {
  getPolicy: (payload: any) => ({ type: PolicyActionTypes.POLICIES_FIND_BY_ID, payload }),
}
interface PolicyParams {
  policyId: string
}
const connector = connect(mapState, mapDispatch)
type EditPolicyProps = ConnectedProps<typeof connector> & RouteComponentProps<PolicyParams>

export class EditPolicy extends Component<EditPolicyProps> {
  static INITIAL_STATE = {
  };
  constructor(props: EditPolicyProps) {
    super(props);
    this.state = {
      ...EditPolicy.INITIAL_STATE,
    };
  }
  componentDidMount() {
    const { match, getPolicy } = this.props
    policiesService.findById(match.params.policyId).then(response => {
      getPolicy(response.data)
    })
  }
  static propKey(propertyName: string, value: any): object {
    return { [propertyName]: value };
  }
  setStateWithEvent(event: any, columnType: string): void {
    this.setState(EditPolicy.propKey(columnType, (event.target as any).value));
  }
  onSubmit = (event: any) => {

  }
  render() {
    const { policy } = this.props
    const header = ['Type', 'Pattern', ''];
    return (
      <Box
        padding="6"
        paddingX={{ base: '4', md: '6' }}
        background="white">
        <form onSubmit={event => this.onSubmit(event)}>
        <SimpleGrid
          display={{ base: "initial", md: "grid" }}
          columns={{ md: 3 }}
          spacing={{ md: 6 }}>
          <GridItem colSpan={{ md: 1 }}>
            <Box px={[4, 0]}>
              <Heading fontSize="lg" fontWeight="medium" lineHeight="6">
                Basic
              </Heading>
            </Box>
          </GridItem>
          <GridItem mt={[5, null, 0]} colSpan={{ md: 2 }}>
            <Stack paddingBottom={5} spacing={4}>
              <FormControl>
                <FormLabel>Enabled</FormLabel>
                <Switch defaultChecked={policy.enabled} colorScheme="brand" onChange={event => this.setStateWithEvent(event, "enabled")} />
              </FormControl>
              <FormControl>
                <FormLabel>Type</FormLabel>
                <Select
                  value={policy.type}
                  isDisabled={true}
                  onChange={event => this.setStateWithEvent(event, "type")}>
                  <option value="security">Security</option>
                </Select>
              </FormControl>
              <FormControl>
                <FormLabel>Display Name</FormLabel>
                <Input type="text" name="display_name" value={policy.display_name} onChange={event => this.setStateWithEvent(event, "display_name")} />
              </FormControl>
              <FormControl>
                <FormLabel>Description</FormLabel>
                <Textarea name="description" value={policy.description} onChange={event => this.setStateWithEvent(event, "description")} />
              </FormControl>
            </Stack>
            <Box
              paddingX={{ base: 4, sm: 6 }}
              paddingY={3}
              background="gray.50"
              textAlign="right">
              <Button
                type="submit"
                colorScheme="brand"
                loadingText="Updating.."
                _focus={{ shadow: "" }}
                fontWeight="md">
                Update
              </Button>
            </Box>
          </GridItem>
          <GridItem colSpan={{ md: 1 }}>
            <Box px={[4, 0]}>
              <Heading fontSize="lg" fontWeight="medium" lineHeight="6">
                Conditions
              </Heading>
            </Box>
          </GridItem>
          <GridItem mt={[5, null, 0]} colSpan={{ md: 2 }}>
            <Table variant="simple">
              <Thead>
                <Tr>
                  {header.map((value) => (
                    <Th key={value}>{value}</Th>
                  ))}
                </Tr>
              </Thead>
              <Tbody>
              {policy.conditions && policy.conditions.map(function(condition){
                return (
                <Tr>
                  <Td>{condition.type}</Td>
                  <Td><Pattern pattern={condition.pattern} editable={false}/></Td>
                  <Td textAlign="right"><IconButton aria-label="Delete" size="sm" icon={<IoTrashOutline />} /></Td>
                </Tr>
                )
              })}
              </Tbody>
            </Table>
            <Box paddingTop={3}>
              <Button
                colorScheme="brand"
                variant="outline"
                fontWeight="md"
                size="sm">
                Add
              </Button>
            </Box>
          </GridItem>
        </SimpleGrid>
        </form>
      </Box>
    )
  }
}

export default withRouter(connector(EditPolicy));
