import Axios, {
  AxiosError,
  AxiosRequestConfig,
  AxiosResponse
} from 'axios'
import { API_PATH } from '../constants'
import parseContentRange, { Range } from './contentRange'
import AuthenticationService from '../services/authentication'
import { clearMessage, setLoading, setMessage, setPagination } from '../store/slices/global'

export interface Condition {
  field: string
  operator: string
  value?: string
}
export interface Sort {
  field: string
  direction: string
}
export interface Query {
  limit: number
  query?: string
  conditions: Condition[]
  offset: number
  sort: Sort[]
}
export interface Result {
  data: any
  error: any,
  pagination: Range
}

const interceptor = (store: any, history: any) => {
  const onRequest = (config: AxiosRequestConfig): AxiosRequestConfig => {
    let authorizationToken = localStorage.getItem("token")
    if (authorizationToken) {
      Axios.defaults.headers.common['Authorization'] = `Bearer ${authorizationToken}`
    }
    store.dispatch(setLoading(true))
    store.dispatch(clearMessage())
    return config
  }
  const onRequestError = (error: AxiosError): Promise<AxiosError> => {
    return Promise.reject(error)
  }
  const onResponse = (response: AxiosResponse): AxiosResponse => {
    store.dispatch(setLoading(false))
    const result = <Result>{
      data: {},
      error: {},
      pagination: {}
    }
    if (response.status >= 200 && response.status < 300) {
      const conntentRange = response.headers["Content-Range".toLowerCase()]
      const range = conntentRange && parseContentRange(conntentRange)
      if (range) {
        result.pagination = range
        store.dispatch(setPagination({
          itemsPerPage: range.end,
          offset: range.start,
          totalItems: range.size
        }))
      }
      result["data"] = response.data
    }
    return response;
  }
  const onResponseError = (error: AxiosError): Promise<AxiosError> => {
    store.dispatch(setLoading(false))
    store.dispatch(setMessage(error?.response?.data))
    if(error.response?.status === 401) {
      history.push('login')
    }
    if(error.response?.status === 403) {
      AuthenticationService.refreshToken()
      history.push('/')
    }
    return Promise.reject(error)
  }
  Axios.defaults.baseURL = API_PATH
  Axios.defaults.headers.common['Accept'] = 'application/json'
  Axios.defaults.headers.common['Content-Type'] = 'application/json'
  Axios.interceptors.request.use(onRequest, onRequestError)
  Axios.interceptors.response.use(onResponse, onResponseError)
}
export default {
  interceptor
}

