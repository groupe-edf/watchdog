import { Table, Thead, Tr, Th, Tbody, useColorModeValue, Icon, SkeletonText, Td, Text, Flex, Badge, Select } from "@chakra-ui/react"
import { Fragment, useEffect, useState } from "react"
import { IoFlashOffOutline, IoMailOutline, IoPersonOutline } from "react-icons/io5"
import { useDispatch, useSelector } from "react-redux"
import { useSearchParams } from "react-router-dom"
import { Commit } from "../../components/Commit"
import { Pagination } from "../../components/Pagination"
import { AppDispatch, RootState } from "../../configureStore"
import { Query } from "../../models"
import { getIssues } from "../../store/slices/issue"

const Issues = () => {
  const header = ['Commit', 'Who', 'Policy', 'Offender', 'Severity']
  const dispatch = useDispatch<AppDispatch>()
  let [searchParams, setSearchParams] = useSearchParams()
  const { categories, pagination } = useSelector((state: RootState) => state.global)
  const { issues } = useSelector((state: RootState) => state.issues)
  const [loading, setLoading] = useState(false)
  const [query, setQuery] = useState<Query>({
    limit: 10
  })
  useEffect(() => {
    setLoading(true)
    dispatch(getIssues({
      limit: Number(searchParams.get('limit')),
      offset: Number(searchParams.get('offset'))
    })).unwrap().then(() => {
      setLoading(false)
    })
  }, [dispatch, searchParams])
  useEffect(() => {
    setSearchParams(query as any)
  }, [query])
  const onChangePage = (data: any) => {
    setQuery({
      conditions: query.conditions,
      limit: data.itemsPerPage,
      offset: data.startIndex,
      sort: query.sort
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
        <Select>
          {categories.filter(category => category.extension === 'handler_type').map((category) => (
            <option value={category.value} key={category.id}>{category.title}</option>
          ))}
        </Select>
      </Flex>
      <Table variant="simple" background={useColorModeValue('white', 'gray.800')}>
        <Thead>
          <Tr>
            {header.map((value) => (
              <Th key={value}>{value}</Th>
            ))}
          </Tr>
        </Thead>
        <Tbody>
          {!loading && issues.length > 0 ? issues.map(function(issue) {
            return (
              <Tr key={issue.id}>
                <Td>
                  <Commit repository={issue.repository?.repository_url} commit={issue.commit}/>
                </Td>
                <Td>
                  <Flex alignItems="center">
                    <Icon as={IoPersonOutline} marginRight={2}/>
                    <Text fontWeight="bold">{issue.commit.author?.name}</Text>
                  </Flex>
                  <Flex alignItems="center">
                    <Icon as={IoMailOutline} marginRight={2}/>
                    <Text>{issue.commit.author?.email}</Text>
                  </Flex>
                </Td>
                <Td>{issue.policy?.display_name}</Td>
                <Td>
                  <Flex alignItems="center">
                    {issue.offender?.object}
                  </Flex>
                  <Flex alignItems="center" maxWidth="200px">
                    <Text>{issue.offender?.value} {issue.offender?.operator} {issue.offender?.operand}</Text>
                  </Flex>
                </Td>
                <Td><Badge variant="outline">{issue.severity}</Badge></Td>
              </Tr>
            )
          }) : (
            (loading ?
              <Tr key="loading">
                <Td colSpan={7}>
                  <SkeletonText noOfLines={4} spacing="4"/>
                </Td>
              </Tr> :
              <Tr key="empty">
                <Td colSpan={7} textAlign="center" paddingX={4}>
                  <Icon fontSize="64" as={IoFlashOffOutline} />
                  <Text marginTop={4}>No gitlab integrations found</Text>
                </Td>
              </Tr>
            )
          )}
        </Tbody>
      </Table>
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

export { Issues }
