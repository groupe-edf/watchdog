import {
  Container,
  Divider
} from '@chakra-ui/react'
import Sidebar from "./sidebar/Sidebar"
import Version from "./Version"
import { useDispatch, useSelector } from "react-redux"
import { getProfile } from "../store/slices/authentication"
import { AppDispatch, RootState } from "../configureStore"
import { Outlet } from 'react-router'


const Layout = () => {
  const dispatch = useDispatch<AppDispatch>()
  const { currentUser } = useSelector((state: RootState) => state.authentication)
  let authorizationToken = localStorage.getItem('token')
  if (!currentUser && authorizationToken) {
    dispatch(getProfile())
  }
  return (
    <Sidebar>
      <Outlet/>
      <Divider marginY={4}/>
      <Container>
        <Version />
      </Container>
    </Sidebar>
  )
}

export { Layout }
