import { StatGroup, Stat, StatLabel, StatNumber, StatHelpText, StatArrow, Box } from '@chakra-ui/react';
import { IssuesList } from './issues';
import { Login } from './authentication';
import { PoliciesList } from './policies';
import { RepositoriesList } from './repositories';
import { RulesList } from './rules';
import { Settings } from './settings';

function Index() {
  return (
    <Box bg="white" borderWidth="1px" borderRadius="lg" p="6">
      <StatGroup>
        <Stat>
          <StatLabel>Issues</StatLabel>
          <StatNumber>55</StatNumber>
        </Stat>
        <Stat>
          <StatLabel>Repositories</StatLabel>
          <StatNumber>2</StatNumber>
        </Stat>
        <Stat>
          <StatLabel>Rules</StatLabel>
          <StatNumber>20</StatNumber>
        </Stat>
        <Stat>
          <StatLabel>Policies</StatLabel>
          <StatNumber>4</StatNumber>
        </Stat>
      </StatGroup>
    </Box>
  )
}

function NotFound() {
  return (
    <Box bg="white" borderWidth="1px" borderRadius="lg" p="6"></Box>
  )
}

export { Index, IssuesList, Login, NotFound, PoliciesList, RepositoriesList, RulesList, Settings }
