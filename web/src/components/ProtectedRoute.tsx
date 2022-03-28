import { Navigate } from 'react-router-dom'
import { Layout } from './Layout'

export type ProtectedRouteProps = {
  isPrivate: boolean
  redirectPath?: string
  children: JSX.Element
}

const ProtectedRoute = ({ isPrivate, redirectPath = '/login', children }: ProtectedRouteProps) => {
  if (isPrivate) {
    if (!localStorage.getItem('user')) {
      return <Navigate to={redirectPath} replace />
    }
    return <Layout/>
  }
  return children
}

export { ProtectedRoute }
