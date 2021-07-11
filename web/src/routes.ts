import React from 'react';
import { IconType } from 'react-icons';
import {
  IoBugOutline,
  IoGitBranchOutline,
  IoHomeOutline,
  IoReceiptOutline,
  IoSettingsOutline,
  IoShieldCheckmarkOutline
} from 'react-icons/io5';
import {
  Index,
  IssuesList,
  Login,
  PoliciesList,
  RepositoriesList,
  RulesList,
  Settings
} from './pages';

export interface Route {
  component: React.ElementType;
  exact?: boolean;
  hide?: boolean;
  icon?: IconType;
  path: string;
  isPrivate?: boolean,
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
    component: IssuesList,
    icon: IoBugOutline,
    isPrivate: true,
    path: '/issues',
    title: 'Issues'
  },
  {
    component: Login,
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
    component: Settings,
    icon: IoSettingsOutline,
    path: '/settings',
    isPrivate: true,
    title: 'Settings'
  }
]

