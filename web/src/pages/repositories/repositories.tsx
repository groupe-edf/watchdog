import {
  Link,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Text,
  useColorModeValue,
  IconButton,
  Spinner,
  Icon,
  Stack,
  Input,
  InputGroup,
  InputLeftElement,
  Flex,
  HStack,
  SkeletonText,
  Badge
} from "@chakra-ui/react"
import {
  IoSearchOutline, IoPlaySharp, IoLockClosedOutline, IoGlobeOutline, IoFlashOffOutline,
} from "react-icons/io5"
import { Component, FC } from "react"
import { connect, ConnectedProps } from "react-redux"
import { Link as ReactRouterLink, Route, RouteComponentProps, Switch, useRouteMatch, withRouter } from "react-router-dom"
import { debounce } from "lodash"
import { withStatusIndicator } from '../../components/withStatusIndicator'
import ShowRepository from './repository'
import repositoryService from "../../services/repository"
import { RepositoryActionTypes, Repository } from "../../store/repositories/types"
import { ApplicationState } from "../../store"
import { StatusBadge } from "../../components/StatusBadge"
import Analyze from "../../components/analyzes/Analyze"
import { Pagination } from "../../components/Pagination"
import { TableState } from "../../store/global/types"

interface RepositoriesContentProps {
  data: Repository[]
  isLoading: boolean
}

export const RepositoriesContent: FC<RepositoriesContentProps> = ({ data, isLoading }) => {
  const header = ['Repository', 'Last Analysis', 'Duration', 'Issues', 'Severity', 'Status', 'Actions'];
  const match = useRouteMatch()
  const analyze = async (event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    event.preventDefault()
    const repository = event.currentTarget.getAttribute("data-repository")
    if (repository != null) {
      repositoryService.analyzeById(repository, {}).then(() => {
        window.location.reload()
      })
    }
  }
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
        {data.length > 0 ? (data.map(function(repository){
          return (
            <Tr key={repository.id}>
              <Td>
                <HStack>
                  {repository.visibility === "public" ? (
                    <Icon as={IoGlobeOutline} />
                  ) : (
                    <Icon as={IoLockClosedOutline} />
                  )}
                  <Link as={ReactRouterLink} color="brand.100" to={`${match.url}/${repository.id}`} style={{ textDecoration: 'none' }}>
                    {repository.repository_url}
                  </Link>
                </HStack>
                <Stack direction="row" marginTop={2}>
                  {repository.integration &&
                  <StatusBadge state={repository.integration.instance_name}></StatusBadge>
                  }
                </Stack>
              </Td>
              <Td>
                {repository.last_analysis?.started_at && new Intl.DateTimeFormat("en-GB", {
                  year: "numeric",
                  month: "long",
                  day: "2-digit",
                  hour: "2-digit",
                  minute: "2-digit",
                  second: "2-digit",
                }).format(Date.parse(repository.last_analysis?.started_at))}
              </Td>
              <Td>{repository.last_analysis?.duration && new Date(repository.last_analysis?.duration / 1000 / 1000).toISOString().substr(11, 8)}</Td>
              <Td>
                <Link as={ReactRouterLink} to={`/issues?conditions=repository_id,eq,${repository.id}`} style={{ textDecoration: 'none' }}>
                  <Text fontWeight="bold" color="brand.100">{repository.last_analysis?.total_issues}</Text>
                </Link>
              </Td>
              <Td><Badge>{repository.last_analysis?.severity}</Badge></Td>
              <Td>
                {repository.last_analysis?.state ? (
                  <StatusBadge state={repository.last_analysis?.state}></StatusBadge>
                ) : (
                  ""
                )}
              </Td>
              <Td align="right">
                {repository.last_analysis?.state != "STARTED" ? (
                  <IconButton data-repository={repository.id} onClick={analyze} aria-label="Run analysis" variant="outline" colorScheme="brand" size="sm" icon={<IoPlaySharp />} />
                ) : (
                  <Spinner />
                )}
              </Td>
            </Tr>
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
                <Text marginTop={4}>No repositories found</Text>
              </Td>
            </Tr>
          )
        ]}
      </Tbody>
    </Table>
  )
}
RepositoriesContent.displayName = 'Repositories'

const mapState = (state: ApplicationState) => ({
  state: state.repositories
})
const mapDispatch = {
  findAll: (payload: any) => ({ type: RepositoryActionTypes.REPOSITORIES_FIND_ALL, payload })
}
const connector = connect(mapState, mapDispatch)
type RepositoriesListProps = ConnectedProps<typeof connector> & RouteComponentProps

class RepositoriesList extends Component<RepositoriesListProps, TableState> {
  constructor(props: RepositoriesListProps) {
    super(props)
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
    this.findAllRepositories()
  }
  componentDidUpdate(prevProps: RepositoriesListProps, prevState: TableState): void {
    if (JSON.stringify(this.state.query) !== JSON.stringify(prevState.query)) {
      this.findAllRepositories()
    }
  }
  findAllRepositories = () => {
    const { findAll } = this.props
    this.setState({ isLoading: true })
    window.scrollTo(0, 0)
    repositoryService.findAll(this.state.query).then(response => {
      findAll(response.data)
      this.setState({ totalItems: response.pagination?.size || 0 })
      this.setState({ isLoading: false })
    })
  }
  static propKey(propertyName: string, value: any): object {
    return { [propertyName]: value }
  }
  handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { query } = { ...this.state }
    const currentQuery = query as any
    const { name, value } = event.target
    currentQuery[name] = value
    this.setState({ query: currentQuery })
  }
  deboucendChangeHandler = (event: React.ChangeEvent<HTMLInputElement>) => {
    const debounced = debounce(() => this.handleChange(event), 300)
    debounced()
  }
  setStateWithEvent(event: any, columnType: string): void {
    this.setState(RepositoriesList.propKey(columnType, (event.target as any).value));
  }
  onChangePage = (data: any) => {
    let query = {...this.state.query}
    query.limit = data.itemsPerPage
    query.offset = data.startIndex
    this.setState({ query: query })
  }
  render() {
    const { match, state } = this.props
    const { isLoading, query, totalItems } = this.state
    return (
      <Switch>
        <Route path={`${match.url}/:repositoryId`} component={ShowRepository}/>
        <Route exact path={match.url}>
          <Flex
            as="header"
            align="center"
            justify="space-between"
            marginBottom={4}
            width="full">
          <InputGroup width="96" display={{ base: "none", md: "flex" }}>
            <InputLeftElement children={<IoSearchOutline/>} />
            <Input
              name="query"
              value={query.query}
              onChange={event => {
                let query = {...this.state.query}
                query.query = event.target.value
                this.setState({ query: query })
              }}
              focusBorderColor="brand.100"
              placeholder="Search for repositories..."
              background="white"/>
          </InputGroup>
          <Flex align="center">
            <Analyze/>
          </Flex>
          </Flex>
          <RepositoriesContent data={state.repositories} isLoading={isLoading}/>
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

export default withRouter(connector(RepositoriesList))
