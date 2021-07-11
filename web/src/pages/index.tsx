import {
  Box,
  Grid,
  GridItem,
  Stat,
  StatHelpText,
  StatLabel,
  StatNumber
} from '@chakra-ui/react';
import { Login, Register } from './authentication'
import { IntegrationsList } from './integrations'
import { IssuesList } from './issues'
import { JobsList } from './jobs'
import { LeaksList } from './leaks'
import { PoliciesList } from './policies'
import Profile from './profile'
import { RepositoriesList } from './repositories'
import RulesList from './rules/rules'
import { Settings } from './settings'
import { UserList } from './users'
import { API_PATH } from '../constants'
import { fetchData } from '../services/commons'
import { Component } from 'react'
import { Label, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid, Pie, PieChart, LabelList, Cell, BarChart, Bar } from 'recharts'
import moment from 'moment'
import { LastAnalyzes } from '../components/analyzes/LastAnalyzes'
import { LeaksBySeverity } from '../components/analytics/LeaksBySeverity';
import { LeaksGraph } from '../components/analytics/LeaksGraph';

export interface DataItem {
  x: string
  y: string
}

class Index extends Component<any, {
  data: any
}> {
  constructor(props: any) {
    super(props)
    this.state = {
      data: [],
    };
  }
  componentDidMount() {
    fetchData<DataItem[]>("GET", `${API_PATH}/analytics`).then(response => this.setState({
      data: response.data
    }))
  }
  render() {
    const { data } = this.state;
    return (
      <>
      <Grid gap={2} templateColumns="repeat(4, 1fr)">
        <GridItem background="white" padding={2} shadow="md">
          <Stat>
            <StatLabel>Repositories</StatLabel>
            <StatNumber>{data.total_items?.repositories}</StatNumber>
            <StatHelpText>Feb 12 - Feb 28</StatHelpText>
          </Stat>
        </GridItem>
        <GridItem background="white" padding={2} shadow="md">
          <Stat>
            <StatLabel>Analyzes</StatLabel>
            <StatNumber>{data.total_items?.repositories_analyzes}</StatNumber>
            <StatHelpText>Feb 12 - Feb 28</StatHelpText>
          </Stat>
        </GridItem>
        <GridItem background="white" padding={2} shadow="md">
          <Stat>
            <StatLabel>Issues</StatLabel>
            <StatNumber>{data.total_items?.repositories_issues}</StatNumber>
            <StatHelpText>Feb 12 - Feb 28</StatHelpText>
          </Stat>
        </GridItem>
        <GridItem background="white" padding={2} shadow="md">
          <Stat>
            <StatLabel>Leaks</StatLabel>
            <StatNumber>{data.total_items?.repositories_leaks}</StatNumber>
            <StatHelpText>Feb 12 - Feb 28</StatHelpText>
          </Stat>
        </GridItem>
        <GridItem height="400px" background="white" padding={2} shadow="md" colSpan={2}>
          <LeaksBySeverity data={data.leak_count_by_severity} />
        </GridItem>
        <GridItem height="400px" background="white" padding={2} shadow="md" colSpan={2}>
          <LeaksGraph data={data.leak_count}/>
        </GridItem>
        <GridItem background="white" padding={2} shadow="md" colSpan={2} display="block">
          <LastAnalyzes width="100%" height="100%"/>
        </GridItem>
      </Grid>
      </>
    )
  }
}

function NotFound() {
  return (
    <Box bg="white" borderWidth="1px" borderRadius="lg" p="6"></Box>
  )
}

export {
  Index,
  IntegrationsList,
  IssuesList,
  JobsList,
  LeaksList,
  Login,
  NotFound,
  PoliciesList,
  Profile,
  RepositoriesList,
  Register,
  RulesList,
  Settings,
  UserList
}
