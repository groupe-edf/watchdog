import { useRoutes } from "react-router-dom"
import { Layout } from "./components/Layout"
import { routes } from "./routes"

const App = () => {
  let protectedRoutes: any = []
  let publicRoutes: any = []
  routes.forEach((route) => {
    if (route.isPrivate) {
      protectedRoutes.push(route)
    } else {
      publicRoutes.push(route)
    }
  })
  publicRoutes = publicRoutes.concat({
    element: <Layout/>,
    children: protectedRoutes
  })
  return useRoutes(publicRoutes)
}

export default App
