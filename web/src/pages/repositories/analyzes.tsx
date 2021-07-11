import {
  Table,
  useColorModeValue,
  Thead,
  Tr,
  Th,
  Tbody,
  Td,
  Text,
  Badge,
  IconButton,
  Link,
  Flex,
  Icon
} from "@chakra-ui/react";
import { Component, FC } from "react";
import { IoGitCommitOutline, IoPerson, IoRocketOutline, IoTimeOutline, IoTrashOutline } from "react-icons/io5";
import { connect, ConnectedProps } from "react-redux";
import { Link as ReactRouterLink, RouteComponentProps, withRouter } from "react-router-dom";
import { StatusBadge } from "../../components/StatusBadge";
import { withStatusIndicator } from "../../components/withStatusIndicator";
import analysisService from "../../services/analysis"
import { ApplicationState } from "../../store";
import { Analysis, Repository, RepositoryActionTypes } from "../../store/repositories/types";

interface AnalyzesContentProps {
  data: Analysis[];
}

export const AnalyzesContent: FC<AnalyzesContentProps> = ({ data }) => {
  const header = ['Started By/At', 'Trigger', 'Duration', 'Severity', 'Total Issues', 'State', ''];
  const deleteById = async (analysisId: string, event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    event.preventDefault();
    analysisService.deleteById(analysisId).then(data => {
    })
  }
  return (
    <Table
      variant="simple"
      background={useColorModeValue('white', 'gray.800')}>
      <Thead>
        <Tr>
          {header.map((value) => (
            <Th key={value}>{value}</Th>
          ))}
        </Tr>
      </Thead>
      <Tbody>
        {data && data.map(function(analysis){
          return (
            <Tr key={analysis.id}>
              <Td>
                <Flex alignItems="center">
                  <Icon as={IoPerson} marginRight={2}/>
                  <Text>
                    {analysis.created_by?.first_name} {analysis.created_by?.last_name}
                  </Text>
                </Flex>
                <Flex alignItems="center">
                  <Icon as={IoTimeOutline} marginRight={2}/>
                  <Text>
                    {analysis.started_at && new Intl.DateTimeFormat("en-GB", {
                      year: "numeric",
                      month: "long",
                      day: "2-digit",
                      hour: "2-digit",
                      minute: "2-digit",
                      second: "2-digit",
                    }).format(Date.parse(analysis.started_at))}
                  </Text>
                </Flex>
                <Flex alignItems="center">
                  <Icon as={IoGitCommitOutline} marginRight={2}/>
                  {analysis.last_commit_hash?.substring(0, 8)}
                </Flex>
              </Td>
              <Td>
                <Flex alignItems="center">
                  <Icon as={IoRocketOutline} marginRight={2}/>
                  <Badge fontWeight="bold">
                    {analysis.trigger}
                  </Badge>
                </Flex>
              </Td>
              <Td>{analysis.duration && new Date(analysis.duration / 1000 / 1000).toISOString().substr(11, 8)}</Td>
              <Td>
                <Badge fontWeight="bold">
                  {analysis.severity}
                </Badge>
              </Td>
              <Td>
                <Link as={ReactRouterLink} to={`/issues?conditions=repository_id,eq,${analysis.repository.id}`} style={{ textDecoration: 'none' }}>
                  <Text fontWeight="bold" color="brand.100">{analysis.total_issues}</Text>
                </Link>
              </Td>
              <Td>
                {analysis.state ? (
                  <StatusBadge state={analysis.state} hint={analysis.state_message} />
                ) : (
                  ""
                )}
              </Td>
              <Td>
                <IconButton onClick={(event) => deleteById(analysis.id, event)} aria-label="Delete" colorScheme="brand" size="sm" icon={<IoTrashOutline />} />
              </Td>
            </Tr>
          )
        })}
      </Tbody>
    </Table>
  )
}
AnalyzesContent.displayName = 'Analyzes';
const AnalyzesWithStatusIndicator = withStatusIndicator(AnalyzesContent);

const mapState = (state: ApplicationState) => ({
  state: state.repositories
})
const mapDispatch = {
  getAnalyzes: (payload: any) => ({ type: RepositoryActionTypes.ANALYZES_FIND_ALL, payload }),
}
const connector = connect(mapState, mapDispatch)
type AnalyzesProps = ConnectedProps<typeof connector> & RouteComponentProps & {
  repository: Repository
}

class AnalyzesList extends Component<AnalyzesProps> {
  constructor(props: AnalyzesProps) {
    super(props);
  }
  componentDidMount() {
    const { repository, state, getAnalyzes } = this.props
    if (state.analyzes.length === 0) {
      analysisService.findAllByRepository(repository.id).then(response => {
        getAnalyzes(response.data)
      })
    }
  }
  render() {
    const { match, state } = this.props
    return (
      <AnalyzesWithStatusIndicator data={state.analyzes}/>
    )
  }
}

export default withRouter(connector(AnalyzesList));
