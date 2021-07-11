import React from 'react';
import { IconType } from 'react-icons';
import {
  IoBugOutline,
  IoGitBranchOutline,
  IoHomeOutline,
  IoListOutline,
  IoReceiptOutline,
  IoSettingsOutline,
  IoShieldCheckmarkOutline
} from 'react-icons/io5';
import {
  Index,
  IntegrationsList,
  IssuesList,
  JobsList,
  LeaksList,
  Login,
  PoliciesList,
  Profile,
  Register,
  RepositoriesList,
  RulesList,
  Settings,
  UserList
} from './pages';
import { Reset } from './pages/authentication/reset';

export interface Route {
  component: React.ElementType;
  exact?: boolean;
  hide?: boolean;
  icon?: IconType;
  path: string;
  isPrivate: boolean,
  title: string;
}

export const routes: Array<Route> = [
  {
    component: Index,
    exact: true,
    icon: IoHomeOutline,
    isPrivate: true,
    path: '/',
    title: 'Dashboard'
  },
  {
    component: RepositoriesList,
    icon: IoGitBranchOutline,
    isPrivate: true,
    path: '/repositories',
    title: 'Repositories'
  },
  {
    component: IntegrationsList,
    hide: true,
    isPrivate: true,
    path: '/integrations',
    title: 'Integrations'
  },
  {
    component: IssuesList,
    icon: IoListOutline,
    isPrivate: true,
    path: '/issues',
    title: 'Issues'
  },
  {
    component: JobsList,
    hide: true,
    icon: IoListOutline,
    isPrivate: true,
    path: '/jobs',
    title: 'Jobs'
  },
  {
    component: LeaksList,
    icon: IoBugOutline,
    isPrivate: true,
    path: '/leaks',
    title: 'Leaks'
  },
  {
    component: Login,
    isPrivate: false,
    hide: true,
    path: '/login',
    title: 'Login'
  },
  {
    component: RulesList,
    icon: IoShieldCheckmarkOutline,
    isPrivate: true,
    path: '/rules',
    title: 'Rules'
  },
  {
    component: PoliciesList,
    icon: IoReceiptOutline,
    path: '/policies',
    isPrivate: true,
    title: 'Policies'
  },
  {
    component: Profile,
    hide: true,
    isPrivate: true,
    path: '/profile',
    title: 'Profile'
  },
  {
    component: Register,
    isPrivate: false,
    hide: true,
    path: '/register',
    title: 'Register'
  },
  {
    component: Reset,
    isPrivate: false,
    hide: true,
    path: '/reset',
    title: 'Reset'
  },
  {
    component: Settings,
    icon: IoSettingsOutline,
    isPrivate: true,
    hide: true,
    path: '/settings',
    title: 'Settings'
  },
  {
    component: UserList,
    isPrivate: true,
    hide: true,
    path: '/users',
    title: 'Users'
  }
]

