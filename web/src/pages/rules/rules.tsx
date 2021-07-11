import { Badge, Box, Button, Stack, Switch, Table, TableCaption, Tbody, Td, Th, Thead, Tr, useColorModeValue, useDisclosure } from "@chakra-ui/react";
import { Component, FC } from "react";
import { IoBuildOutline } from "react-icons/io5";
import { connect, ConnectedProps } from "react-redux";
import { RouteComponentProps, withRouter } from "react-router-dom";
import { ButtonLink } from "../../components/ButtonLink";
import { Pagination } from "../../components/Pagination";
import { withStatusIndicator } from "../../components/withStatusIndicator";
import policy from "../../services/policy";
import rulesService from "../../services/rule"
import { ApplicationState } from "../../store";
import { Rule, RuleActionTypes } from "../../store/rules/types";

interface RulesContentProps {
  data: Rule[];
}

export const RulesContent: FC<RulesContentProps> = ({ data }) => {
  const header = ['Name', 'Severity', 'Tags', 'State', ''];
  const toggleRule = (ruleId: number, enabled: boolean, event: React.ChangeEvent<HTMLInputElement>) => {
    rulesService.toggle(ruleId, !enabled).then(data => {
    })
  };
  return (
    <>
      <Box paddingBottom={4}>
        <Stack
          justify={'flex-end'}
          direction={'row'}>
          <Button
            colorScheme="brand"
            fontWeight="md">
            New
          </Button>
        </Stack>
      </Box>
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
          {data && data.map(function(rule){
            return (
              <Tr key={rule.id}>
                <Td>{rule.display_name}</Td>
                <Td><Badge>{rule.severity}</Badge></Td>
                <Td>
                  <Stack direction="row">
                    {rule.tags && rule.tags.map(function(tag){
                      return (
                        <Badge variant="outline" colorScheme="brand" key={tag}>
                          {tag}
                        </Badge>
                      )
                    })}
                  </Stack>
                </Td>
                <Td><Switch defaultChecked={rule.enabled} colorScheme="brand" onChange={(event) => toggleRule(rule.id, rule.enabled, event)}/></Td>
                <Td textAlign="right">
                  <ButtonLink rightIcon={<IoBuildOutline />} to={`/rules/${rule.id}/edit`} colorScheme="gray" size="sm" variant="outline">
                    Edit
                  </ButtonLink>
                </Td>
              </Tr>
            )
          })}
        </Tbody>
      </Table>
    </>
  )
}
RulesContent.displayName = 'Rules';
const RulesWithStatusIndicator = withStatusIndicator(RulesContent);

const mapState = (state: ApplicationState) => ({
  state: state.rules
})
const mapDispatch = {
  getRules: (payload: any) => ({ type: RuleActionTypes.RULES_FIND_ALL, payload }),
}
const connector = connect(mapState, mapDispatch)
type RulesProps = ConnectedProps<typeof connector> & RouteComponentProps

class RulesList extends Component<RulesProps> {
  constructor(props: RulesProps) {
    super(props);
  }
  componentDidMount() {
    const { state, getRules } = this.props
    rulesService.findAll().then(response => {
      getRules(response.data)
    })
  }
  render() {
    const { match, state } = this.props
    return (
      <RulesWithStatusIndicator data={state.rules}/>
    )
  }
}

export default withRouter(connector(RulesList));
