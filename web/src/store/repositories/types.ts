import { Action } from "redux";
import { Integration } from "../../store/integrations/types"
import { User } from "../users/types";

export enum RepositoryActionTypes {
  ANALYZE = "@@analyzes/ANALYZE",
  ANALYZES_FIND_ALL = "@@analyzes/FIND_ALL",
  REPOSITORIES_APPEND = "@@repositories/APPEND",
  REPOSITORIES_APPEND_ALL = "@@repositories/APPEND_ALL",
  REPOSITORIES_DELETE_ALL = "@@repositories/DELETE_ALL",
  REPOSITORIES_DELETE_BY_ID = "@@repositories/DELETE_BY_ID",
  REPOSITORIES_FIND_ALL = "@@repositories/FIND_ALL",
  REPOSITORIES_FIND_BY_ID = "@@repositories/FIND_BY_ID"
}

export interface Repository {
  id: string
  integration?: Integration
  last_analysis?: Analysis
  repository_url: string
  issues?: number
  visibility: string
}

export interface Analysis {
  id: string
  created_by?: User
  duration?: number
  finished_at?: string
  last_commit_hash: string
  repository: Repository
  started_at?: string
  state: string
  state_message?: string
  severity: string
  total_issues: number
  trigger: string
}

export interface RepositoryState {
  readonly analyzes: Analysis[],
  readonly repositories: Repository[],
  readonly repository: Repository;
}

export interface AppendAction extends Action {
  type: RepositoryActionTypes.REPOSITORIES_APPEND;
  payload: {
    state: RepositoryState;
  };
}

export interface AnalyzesFindAllAction extends Action {
  type: RepositoryActionTypes.ANALYZES_FIND_ALL;
  payload: {
    state: RepositoryState;
  };
}

export interface RepositoriesDeleteByIdAction extends Action {
  type: RepositoryActionTypes.REPOSITORIES_DELETE_BY_ID;
  payload: {
    state: RepositoryState;
  };
}

export interface RepositoriesFindAllAction extends Action {
  type: RepositoryActionTypes.REPOSITORIES_FIND_ALL;
  payload: {
    state: RepositoryState;
  };
}

export interface RepositoriesFindByIdAction extends Action {
  type: RepositoryActionTypes.REPOSITORIES_FIND_BY_ID;
  payload: {
    state: RepositoryState;
  };
}

export type RepositoriesActions =
  AppendAction |
  AnalyzesFindAllAction |
  RepositoriesDeleteByIdAction |
  RepositoriesFindAllAction |
  RepositoriesFindByIdAction;
