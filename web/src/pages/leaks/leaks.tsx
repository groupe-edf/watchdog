import { ExternalLinkIcon } from "@chakra-ui/icons";
import { Icon, Table, Text, Tbody, Td, Th, Thead, Tr, useColorModeValue, Link, MenuButton, Menu, MenuItem, MenuList, Checkbox, Badge, Stack, Box, Select, FormControl, FormLabel, HStack, Popover, PopoverArrow, PopoverBody, PopoverCloseButton, PopoverContent, PopoverHeader, PopoverTrigger, IconButton, Code, SkeletonText, LinkBox, LinkOverlay } from "@chakra-ui/react";
import { Component, FC } from "react";
import { IoEllipsisVertical, IoEyeOutline, IoFlashOffOutline } from "react-icons/io5";
import { connect, ConnectedProps } from "react-redux";
import { Link as RouterLink, Route, RouteComponentProps, Switch, withRouter } from "react-router-dom";
import { Pagination } from "../../components/Pagination";
import { withStatusIndicator } from "../../components/withStatusIndicator";
import leaksService from "../../services/leak"
import rulesService from "../../services/rule"
import { ApplicationState } from "../../store";
import { TableState } from "../../store/global/types";
import { Leak, LeakActionTypes } from "../../store/leaks/types";
import { RuleActionTypes } from "../../store/rules/types";
import ShowLeak from "./leak";

interface LeaksContentProps {
  data: Leak[]
  isLoading: boolean
}

export const LeaksContent: FC<LeaksContentProps> = ({ data, isLoading }) => {
  const header = ['#', 'Rule', 'Severity', 'File', 'Author', 'Offender', '']
  return (
    <>
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
          {!isLoading && data.length > 0 ? (data.map(function(leak, index){
            return (
              <LinkBox as="tr" key={index}>
                <Td><Checkbox></Checkbox></Td>
                <Td>
                  <LinkOverlay as={RouterLink} to={`/leaks/${leak.id}`} style={{ textDecoration: 'none' }}>
                    {leak.rule.display_name}
                  </LinkOverlay>
                  <Stack direction="row">
                    {leak.rule.tags && leak.rule.tags.map(function(tag){
                      return (
                        <Badge variant="outline" colorScheme="brand" key={tag}>
                          {tag}
                        </Badge>
                      )
                    })}
                  </Stack>
                </Td>
                <Td><Badge>{leak.severity}</Badge></Td>
                <Td>
                  <Link href={leak.repository.repository_url + "/commit/" + leak.commit_hash} color="brand.100" isExternal>
                    {leak.repository.repository_url} <ExternalLinkIcon mx="2px" />
                  </Link>
                  <Text>{leak.file}</Text>
                  <Text>Line: {leak.line_number}</Text>
                </Td>
                <Td>
                  <Text fontWeight="bold">{leak.author}</Text>
                  <Text>{leak.created_at && new Intl.DateTimeFormat("en-GB", {
                    year: "numeric",
                    month: "long",
                    day: "2-digit",
                    hour: "2-digit",
                    minute: "2-digit",
                    second: "2-digit",
                  }).format(Date.parse(leak.created_at))}</Text>
                </Td>
                <Td>
                  <Popover placement="bottom-start">
                    <PopoverTrigger>
                      <IconButton aria-label="Reveal" colorScheme="gray" size="sm" icon={<IoEyeOutline />} />
                    </PopoverTrigger>
                    <PopoverContent>
                      <PopoverArrow />
                      <PopoverCloseButton />
                      <PopoverHeader>Offender</PopoverHeader>
                      <PopoverBody>
                        <Text>{leak.offender}</Text>
                        <Code>{leak.line}</Code>
                      </PopoverBody>
                    </PopoverContent>
                  </Popover>
                </Td>
                <Td>
                  <Menu>
                    <MenuButton>
                    <IoEllipsisVertical />
                    </MenuButton>
                    <MenuList>
                      <MenuItem>Mark as false positive</MenuItem>
                      <MenuItem>Mark as resolved</MenuItem>
                    </MenuList>
                  </Menu>
                </Td>
              </LinkBox>
            )
          })) : [
            (isLoading ?
              <Tr key="loading">
                <Td colSpan={7}>
                  <SkeletonText noOfLines={4} spacing="4" />
                </Td>
              </Tr> :
              <Tr key="empty">
                <Td colSpan={7} textAlign="center" color="grey" paddingX={4}>
                  <Icon fontSize="64" as={IoFlashOffOutline} />
                  <Text marginTop={4}>No leaks found</Text>
                </Td>
              </Tr>
            )
          ]}
        </Tbody>
      </Table>
    </>
  )
}
LeaksContent.displayName = 'Leaks';
const LeaksWithStatusIndicator = withStatusIndicator(LeaksContent);

const mapState = (state: ApplicationState) => ({
  state: state
})
const mapDispatch = {
  getLeaks: (payload: any) => ({ type: LeakActionTypes.LEAKS_FIND_ALL, payload }),
  getRules: (payload: any) => ({ type: RuleActionTypes.RULES_FIND_ALL, payload }),
}
const connector = connect(mapState, mapDispatch)
type LeaksProps = ConnectedProps<typeof connector> & RouteComponentProps

class LeaksList extends Component<LeaksProps, TableState> {
  constructor(props: LeaksProps) {
    super(props);
    this.state = {
      isLoading: false,
      query: {
        conditions: [],
        limit: 10,
        offset: 0,
        sort: []
      },
      totalItems: 0
    }
  }
  componentDidMount() {
    const { state, getRules } = this.props
    this.findAllLeaks()
    rulesService.findAll().then(response => {
      getRules(response.data)
    })
  }
  componentDidUpdate(prevProps: LeaksProps, prevState: TableState): void {
    if (!Object.is(this.state.query, prevState.query)) {
      this.findAllLeaks()
    }
  }
  findAllLeaks = () => {
    const { state, getLeaks } = this.props
    this.setState({ isLoading: true })
    window.scrollTo(0, 0)
    leaksService.findAll(this.state.query).then(response => {
      getLeaks(response.data)
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
    const { match, state } = this.props
    const { isLoading, query, totalItems } = this.state;
    return (
      <Switch>
        <Route path={`${match.url}/:leakId`} component={ShowLeak}/>
        <Route exact path={match.url}>
          <Box padding={4} background="#fafafa" borderBottom={1} borderBottomColor="#dbdbdb">
            <HStack isInline={true}>
              <FormControl id="rule">
                <FormLabel>Rule</FormLabel>
                <Select
                  backgroundColor="white"
                  onChange={event => {
                    let query = {...this.state.query}
                    query.conditions = [
                      { field: "rule_id", operator: "eq", value: event.target.value }
                    ]
                    this.setState({ query: query })
                  }}
                  placeholder="Select rule">
                  {state.rules.rules.length > 0 && (state.rules.rules.map(function(rule){
                    return (
                      <option key={rule.id} value={rule.id}>{rule.display_name}</option>
                    )
                  }))}
                </Select>
              </FormControl>
              <FormControl id="severity">
                <FormLabel>Severity</FormLabel>
                <Select
                  backgroundColor="white"
                  onChange={event => {
                    let query = {...this.state.query}
                    query.conditions = [
                      { field: "severity", operator: "eq", value: event.target.value }
                    ]
                    this.setState({ query: query })
                  }}
                  placeholder="Select severity">
                </Select>
              </FormControl>
            </HStack>
          </Box>
          <LeaksContent data={state.leaks.leaks} isLoading={isLoading}/>
          <Pagination
            currentPage={1}
            pagesToShow={5}
            itemsPerPage={query.limit}
            offset={query.offset}
            onChangePage={this.onChangePage}
            totalItems={totalItems}/>
        </Route>
      </Switch>
    )
  }
}

export default withRouter(connector(LeaksList))
