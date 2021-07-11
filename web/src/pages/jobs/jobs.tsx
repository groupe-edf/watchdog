import { Component, FC } from "react"
import { connect, ConnectedProps } from "react-redux"
import { RouteComponentProps, withRouter } from "react-router-dom"
import { ApplicationState } from "../../store"
import { TableState } from "../../store/global/types"
import { Job, JobActionTypes } from "../../store/jobs/types"
import jobsService from "../../services/job"
import { useColorModeValue } from "@chakra-ui/color-mode"
import { Table, Thead, Tr, Th, Icon, SkeletonText, Tbody, Td, Text } from "@chakra-ui/react"
import { IoFlashOffOutline } from "react-icons/io5"

interface JobsContentProps {
  data: Job[]
  isLoading: boolean
}

export const JobsContent: FC<JobsContentProps> = ({ data, isLoading }) => {
  const header = ['#', 'Type', 'Started At', 'Queue', 'Priority', 'Error Count', '']
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
        {!isLoading && data.length > 0 ? (data.map(function(job, index){
          return (
            <Tr>
              <Td>{job.id}</Td>
              <Td>{job.type}</Td>
              <Td>
                <Text>{job.started_at && new Intl.DateTimeFormat("en-GB", {
                  year: "numeric",
                  month: "long",
                  day: "2-digit",
                  hour: "2-digit",
                  minute: "2-digit",
                  second: "2-digit",
                }).format(Date.parse(job.started_at))}</Text>
              </Td>
              <Td>{job.queue}</Td>
              <Td>{job.priority}</Td>
              <Td>{job.error_count}</Td>
              <Td></Td>
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
                <Text marginTop={4}>No leaks found</Text>
              </Td>
            </Tr>
          )
        ]}
      </Tbody>
    </Table>
  )
}

const mapState = (state: ApplicationState) => ({
  state: state
})
const mapDispatch = {
  getJobs: (payload: any) => ({ type: JobActionTypes.JOBS_FIND_ALL, payload }),
}
const connector = connect(mapState, mapDispatch)
type JobsProps = ConnectedProps<typeof connector> & RouteComponentProps

class JobsList extends Component<JobsProps, TableState> {
  constructor(props: JobsProps) {
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
    const { state, getJobs } = this.props
    jobsService.findAll().then(response => {
      getJobs(response.data)
    })
  }
  render() {
    const { state } = this.props
    const { isLoading, query, totalItems } = this.state;
    return (
      <JobsContent data={state.jobs.jobs} isLoading={isLoading}/>
    )
  }
}

export default withRouter(connector(JobsList));
