import {
  Box,
  Text
} from '@chakra-ui/react'
import { Login, Register } from './authentication'
import { Repositories } from './repositories'
import { Dashboard } from './Dashboard'
import { Settings } from './settings'
import { Users } from './users/Users'
import { Policies } from './policies/Policies'
import { EditPolicy } from './policies/EditPolicy'
import { Rules } from './rules/Rules'
import { Issues } from './issues/Issues'
import { IntegrationView } from './integrations/IntegrationView'
import { Integrations } from './integrations/Integrations'
import RepositoryView from './repositories/RepositoryView'
import Profile from './authentication/Profile'
import { Leaks, LeakView } from './leaks'
import { AnalysisView, Analyzes } from './analyzes'

function NotFound() {
  return (
    <Box bg="white" borderWidth="1px" borderRadius="lg" p="6">
      <Text>There's nothing here: 404!</Text>
    </Box>
  )
}
export {
  Analyzes,
  AnalysisView,
  Dashboard,
  Integrations,
  IntegrationView,
  Issues,
  Leaks,
  LeakView,
  Login,
  NotFound,
  Policies,
  Profile,
  EditPolicy,
  Register,
  Repositories,
  RepositoryView,
  Rules,
  Settings,
  Users
}
