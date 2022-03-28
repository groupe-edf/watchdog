import { Badge, Flex, HStack, Icon, Input, InputGroup, InputLeftElement, Link, SkeletonText, Stack, Table, Tbody, Td, Text, Th, Thead, Tr, useColorModeValue } from "@chakra-ui/react"
import { Fragment, useEffect, useState } from "react"
import { Link as ReactRouterLink, useSearchParams } from "react-router-dom"
import { IoFlashOffOutline, IoGlobeOutline, IoLockClosedOutline, IoSearchOutline } from "react-icons/io5"
import { useDispatch, useSelector } from "react-redux"
import { StatusBadge } from "../../components/StatusBadge"
import { AppDispatch, RootState } from "../../configureStore"
import { Repository } from "../../models"
import { getRepositories } from "../../store/slices/repository"
import Analyze from "./Analyze"
import { Pagination } from "../../components/Pagination"

export type RepositoryId = {
  repositoryId: string
}
const Items = ({ data, isLoading }: any) => {
  const header = ['Repository', 'Last Analysis', 'Duration', 'Issues', 'Severity', 'Status', 'Actions']
  return (
    <Table variant="simple" background={useColorModeValue('white', 'gray.800')}>
      <Thead>
        <Tr>
          {header.map((value) => (
            <Th key={value}>{value}</Th>
          ))}
        </Tr>
      </Thead>
      <Tbody>
        {data.length > 0 ? (data.map(function(repository: Repository){
          return (
            <Tr key={repository.id}>
              <Td>
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
                <Stack direction="row" marginTop={2}>
                  {repository.integration &&
                  <StatusBadge state={repository.integration.instance_name}></StatusBadge>
                  }
                </Stack>
              </Td>
              <Td>
                {repository.last_analysis?.started_at && new Intl.DateTimeFormat("en-GB", {
                  year: "numeric",
                  month: "long",
                  day: "2-digit",
                  hour: "2-digit",
                  minute: "2-digit",
                  second: "2-digit",
                }).format(Date.parse(repository.last_analysis?.started_at))}
              </Td>
              <Td>{repository.last_analysis?.duration && new Date(repository.last_analysis?.duration / 1000 / 1000).toISOString().substr(11, 8)}</Td>
              <Td>
                <Link as={ReactRouterLink} to={`/issues?conditions=repository_id,eq,${repository.id}`} style={{ textDecoration: 'none' }}>
                  <Text fontWeight="bold" color="brand.100">{repository.last_analysis?.total_issues}</Text>
                </Link>
              </Td>
              <Td><Badge>{repository.last_analysis?.severity}</Badge></Td>
              <Td>
                {repository.last_analysis?.state ? (
                  <StatusBadge state={repository.last_analysis?.state}></StatusBadge>
                ) : (
                  ""
                )}
              </Td>
              <Td align="right">
                <Analyze repository={repository} size='sm'/>
              </Td>
            </Tr>
          )
        })) : [
          (isLoading ?
            <Tr key="loading">
              <Td colSpan={7}>
                <SkeletonText noOfLines={4} spacing="4" />
              </Td>
            </Tr> :
            <Tr key="empty">
              <Td colSpan={7} textAlign="center" color="grey" paddingX={4}>
                <Icon fontSize="64" as={IoFlashOffOutline} />
                <Text marginTop={4}>No repositories found</Text>
              </Td>
            </Tr>
          )
        ]}
      </Tbody>
    </Table>
  )
}

const Repositories = () => {
  const dispatch = useDispatch<AppDispatch>()
  let [searchParams, setSearchParams] = useSearchParams()
  const { pagination } = useSelector((state: RootState) => state.global)
  const { repositories } = useSelector((state: RootState) => state.repositories)
  const [loading, setLoading] = useState(false)
  useEffect(() => {
    setLoading(true)
    dispatch(getRepositories({
      limit: Number(searchParams.get('limit')),
      offset: Number(searchParams.get('offset'))
    })).unwrap().then(() => {
      setLoading(false)
    })
  }, [dispatch, searchParams])
  const onChangePage = (data: any) => {
    setSearchParams({
      limit: data.itemsPerPage,
      offset: data.startIndex
    })
  }
  return (
    <Fragment>
      <Flex
        as="header"
        align="center"
        justify="space-between"
        marginBottom={4}
        width="full">
        <InputGroup width="96" display={{ base: "none", md: "flex" }}>
          <InputLeftElement children={<IoSearchOutline/>} />
          <Input
            name="query"
            onChange={event => setSearchParams({ conditions: `repository_url=${event.target.value}` })}
            placeholder="Search for repositories..."/>
        </InputGroup>
        <Analyze/>
      </Flex>
      <Items data={repositories} isLoading={loading}/>
      <Pagination
        currentPage={1}
        pagesToShow={5}
        itemsPerPage={pagination.itemsPerPage}
        offset={pagination.offset}
        onChangePage={onChangePage}
        totalItems={pagination.totalItems}/>
    </Fragment>
  )
}

export default Repositories
