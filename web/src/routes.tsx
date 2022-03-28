import { IconType } from 'react-icons'
import {
  IoBugOutline,
  IoExtensionPuzzleOutline,
  IoGitBranchOutline,
  IoHomeOutline,
  IoListOutline,
  IoPeopleOutline,
  IoReceiptOutline,
  IoSettingsOutline,
  IoShieldCheckmarkOutline
} from 'react-icons/io5'
import { Navigate, RouteObject } from 'react-router-dom'
import {
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
} from './pages'

export interface Route extends RouteObject {
  hide?: boolean
  icon?: IconType
  isPrivate?: boolean
  title?: string
}

export const routes: Array<Route> = [
  {
    element: <Dashboard/>,
    icon: IoHomeOutline,
    index: true,
    isPrivate: true,
    path: '/',
    title: 'Dashboard'
  }, {
    children: [
      {
        element: <Analyzes/>,
        index: true,
      }
    ],
    element: <Analyzes/>,
    icon: IoListOutline,
    isPrivate: true,
    path: '/analyzes',
    title: 'Analyzes'
  }, {
    children: [
      {
        element: <Integrations/>,
        index: true,
      }, {
        element: <IntegrationView/>,
        path: ':integration_id'
      }
    ],
    icon: IoExtensionPuzzleOutline,
    isPrivate: true,
    hide: false,
    path: 'integrations',
    title: 'Integrations'
  }, {
    children: [
      {
        element: <Leaks/>,
        index: true,
      }, {
        element: <LeakView/>,
        path: ':leak_id'
      }
    ],
    icon: IoBugOutline,
    isPrivate: true,
    hide: false,
    path: 'leaks',
    title: 'Leaks'
  }, {
    element: <Issues/>,
    icon: IoListOutline,
    isPrivate: true,
    path: '/issues',
    title: 'Issues'
  }, {
    element: <Login/>,
    isPrivate: false,
    hide: true,
    path: 'login',
    title: 'Login'
  }, {
    children: [
      {
        element: <Policies/>,
        index: true,
      }, {
        element: <EditPolicy/>,
        path: ':policy_id/edit'
      }
    ],
    icon: IoReceiptOutline,
    isPrivate: true,
    path: '/policies',
    title: 'Policies'
  }, {
    element: <Profile/>,
    hide: true,
    isPrivate: true,
    path: 'profile',
    title: 'Profile'
  }, {
    children: [
      {
        element: <Repositories/>,
        index: true,
      }, {
        element: <RepositoryView/>,
        path: ':repository_id'
      }
    ],
    icon: IoGitBranchOutline,
    isPrivate: true,
    path: 'repositories',
    title: 'Repositories'
  }, {
    element: <Register/>,
    hide: true,
    path: 'register'
  }, {
    element: <Rules/>,
    icon: IoShieldCheckmarkOutline,
    isPrivate: true,
    path: '/rules',
    title: 'Rules'
  }, {
    element: <Settings/>,
    icon: IoSettingsOutline,
    isPrivate: true,
    path: "settings",
    title: 'Settings'
  }, {
    element: <Navigate replace to="/"/>,
    isPrivate: false,
    hide: true,
    path: 'redirect',
    title: 'Redirect'
  }, {
    element: <Users/>,
    icon: IoPeopleOutline,
    isPrivate: true,
    path: 'users',
    title: 'Users'
  }
]
