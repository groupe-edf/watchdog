import { Table, useColorModeValue, Thead, Tr, Th, Tbody, Td, Box, Stack, Select, Icon, Text, Flex, Badge, SkeletonText } from "@chakra-ui/react";
import { Component, FC } from "react";
import { IoFlashOffOutline, IoMailOutline, IoPersonOutline } from "react-icons/io5";
import { connect, ConnectedProps } from "react-redux";
import { RouteComponentProps, withRouter } from "react-router-dom";
import { Commit } from "../../components/Commit";
import { Pagination } from "../../components/Pagination";
import { withStatusIndicator } from "../../components/withStatusIndicator";
import issuesService from "../../services/issue"
import policiesService from "../../services/policy"
import { ApplicationState } from "../../store";
import { TableState } from "../../store/global/types";
import { Issue, IssueActionTypes } from "../../store/issues/types";
import { PolicyActionTypes } from "../../store/policies/types";

interface IssuesContentProps {
  data: Issue[]
  isLoading: boolean
}

export const IssuesContent: FC<IssuesContentProps> = ({ data, isLoading }) => {
  const header = ['Commit', 'Who', 'Policy', 'Offender', 'Severity'];
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
        {!isLoading && data.length > 0 ? data.map(function(issue){
          return (
            <Tr key={issue.id}>
              <Td>
                <Commit repository={issue.repository.repository_url} commit={issue.commit}/>
              </Td>
              <Td>
                <Flex alignItems="center">
                  <Icon as={IoPersonOutline} marginRight={2}/>
                  <Text fontWeight="bold">{issue.commit.author}</Text>
                </Flex>
                <Flex alignItems="center">
                  <Icon as={IoMailOutline} marginRight={2}/>
                  <Text>{issue.commit.email}</Text>
                </Flex>
              </Td>
              <Td>{issue.policy.display_name}</Td>
              <Td>
                <Flex alignItems="center">
                  {issue.offender?.object}
                </Flex>
                <Flex alignItems="center" maxWidth="200px">
                  <Text isTruncated>{issue.offender?.value} {issue.offender?.operator} {issue.offender?.operand}</Text>
                </Flex>
              </Td>
              <Td><Badge variant="outline">{issue.severity}</Badge></Td>
            </Tr>
          )
        }) : [
          (isLoading ?
            <Tr key="loading">
              <Td colSpan={7}>
                <SkeletonText noOfLines={4} spacing="4" />
              </Td>
            </Tr> :
            <Tr key="empty">
              <Td colSpan={7} textAlign="center" color="grey" paddingX={4}>
                <Icon fontSize="64" as={IoFlashOffOutline} />
                <Text marginTop={4}>No issues found</Text>
              </Td>
            </Tr>
          )
        ]}
      </Tbody>
    </Table>
  )
}
IssuesContent.displayName = 'Issues';
const IssuesWithStatusIndicator = withStatusIndicator(IssuesContent);

const mapState = (state: ApplicationState) => ({
  state: state
})
const mapDispatch = {
  getIssues: (payload: any) => ({ type: IssueActionTypes.ISSUES_FIND_ALL, payload }),
  getPolicies: (payload: any) => ({ type: PolicyActionTypes.POLICIES_FIND_ALL, payload }),
}
const connector = connect(mapState, mapDispatch)
type IssuesProps = ConnectedProps<typeof connector> & RouteComponentProps

class IssuesList extends Component<IssuesProps, TableState> {
  constructor(props: IssuesProps) {
    super(props);
    this.findAllIssues = this.findAllIssues.bind(this);
    this.state = {
      isLoading: false,
      query: {
        conditions: [],
        limit: 10,
        offset: 0,
        sort: []
      },
      totalItems: 0,
    }
  }
  componentDidMount() {
    const { state, getPolicies } = this.props
    this.findAllIssues()
    if (state.policies.policies.length === 0) {
      policiesService.findAll().then(response => {
        getPolicies(response.data)
      })
    }
  }
  componentDidUpdate(prevProps: IssuesProps, prevState: TableState): void {
    if (JSON.stringify(this.state.query) !== JSON.stringify(prevState.query)) {
      console.log(JSON.stringify(prevState.query))
      console.log(JSON.stringify(this.state.query))
      this.findAllIssues()
    }
  }
  findAllIssues = () => {
    const { getIssues } = this.props
    this.setState({ isLoading: true })
    window.scrollTo(0, 0)
    issuesService.findAll(this.state.query).then(response => {
      getIssues(response.data)
      this.setState({ totalItems: response.pagination?.size || 0 })
      this.setState({ isLoading: false })
    })
  }
  onChangePage = (data: any) => {
    let query = {...this.state.query}
    query.limit = data.itemsPerPage
    query.offset = data.startIndex
    this.setState({ query: query })
  }
  render() {
    const { state } = this.props
    const { isLoading, query, totalItems } = this.state
    return (
      <>
        <Box paddingBottom={4}>
          <Stack
            justify="space-between"
            direction="row"
            isInline={true}>
            <Select backgroundColor="white"
              onChange={event => {
                let query = {...this.state.query}
                query.conditions = [
                  { field: "policy_id", operator: "eq", value: event.target.value }
                ]
                this.setState({ query: query })
              }}
              placeholder="Select policy">
              {state.policies.policies.length > 0 && (state.policies.policies.map(function(policy){
                return (
                  <option key={policy.id} value={policy.id}>{policy.display_name}</option>
                )
              }))}
            </Select>
            <Select backgroundColor="white"
              onChange={event => {
                let query = {...this.state.query}
                query.conditions = [
                  { field: "severity", operator: "eq", value: event.target.value }
                ]
                this.setState({ query: query })
              }}
              placeholder="Select severity">
              {state.categories.categories.length > 0 && (state.categories.categories.map(function(category){
                if (category.extension === "issue_severity") {
                  return (
                    <option key={category.id} value={category.value || category.id}>{category.title}</option>
                  )
                }
              }))}
            </Select>
            <Select value={query.limit} backgroundColor="white"
              onChange={event => {
                let query = {...this.state.query}
                query.limit = parseInt(event.target.value)
                this.setState({ query: query })
              }}>
              <option value={5}>5</option>
              <option value={10}>10</option>
              <option value={25}>25</option>
              <option value={50}>50</option>
              <option value={100}>100</option>
            </Select>
          </Stack>
        </Box>
        <IssuesContent data={state.issues.issues} isLoading={isLoading}/>
        <Pagination
          currentPage={1}
          pagesToShow={5}
          itemsPerPage={query.limit}
          offset={query.offset}
          onChangePage={this.onChangePage}
          totalItems={totalItems}/>
      </>
    )
  }
}

export default withRouter(connector(IssuesList));
