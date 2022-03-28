import { Grid, GridItem, Heading, HStack, Icon, Link, Stack } from "@chakra-ui/react"
import { Card, CardBody, CardHeader, CardTitle } from "@saas-ui/react"
import { Fragment, useEffect, useState } from "react"
import { IoGlobeOutline, IoLockClosedOutline } from "react-icons/io5"
import { useDispatch, useSelector } from "react-redux"
import { Link as ReactRouterLink, useParams } from "react-router-dom"
import { AppDispatch, RootState } from "../../configureStore"
import { getAnalysisById } from "../../store/slices/analysis"
import { getRepositoryById } from "../../store/slices/repository"
import AnalysisView from "../analyzes/AnalysisView"
import Analyze from "./Analyze"
import DeleteRepository from "./DeleteRepository"
import RepositoryBadge from "./RepositoryBadge"

export type RepositoryId = {
  repository_id: string
}
const RepositoryView = () => {
  const dispatch = useDispatch<AppDispatch>()
  const { repository_id } = useParams<RepositoryId>() as RepositoryId
  const { repository } = useSelector((state: RootState) => state.repositories)
  const { analysis } = useSelector((state: RootState) => state.analyzes)
  const [loading, setLoading] = useState(false)
  useEffect(() => {
    setLoading(true)
    dispatch(getRepositoryById(repository_id)).unwrap().then((response) => {
      if (response.last_analysis) {
        dispatch(getAnalysisById(response.last_analysis.id))
      }
      setLoading(false)
    })
  }, [])
  return (
    <Fragment>
      <Stack spacing={4} direction="row-reverse" align="center" marginBottom={4}>
        <Analyze repository={repository}/>
        <DeleteRepository repositoryId={repository.id}/>
      </Stack>
      <Grid gap={4} templateColumns='repeat(6, 1fr)'>
        <GridItem colSpan={4}>
          <Card>
            <CardHeader>
              <CardTitle fontSize="xl">Details</CardTitle>
            </CardHeader>
            <CardBody>
              {!loading && repository.last_analysis && <RepositoryBadge repositoryId={repository.id}/>}
              <HStack>
                {repository.visibility === "public" ? (
                  <Icon as={IoGlobeOutline} />
                ) : (
                  <Icon as={IoLockClosedOutline} />
                )}
                <Link as={ReactRouterLink} color="brand.100" to={repository.id} style={{ textDecoration: 'none' }}>
                  {repository.repository_url}
                </Link>
              </HStack>
            </CardBody>
          </Card>
        </GridItem>
        <GridItem colSpan={2}>
          <Card>
            <CardHeader>
              <CardTitle>Last Analysis</CardTitle>
            </CardHeader>
            <CardBody>
              {analysis && <AnalysisView analysis={analysis}/>}
            </CardBody>
          </Card>
        </GridItem>
      </Grid>
    </Fragment>
  )
}

export default RepositoryView
