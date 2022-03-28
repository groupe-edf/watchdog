import { Icon, Link, SkeletonText, Table, Tbody, Td, Th, Thead, Tr, Text, useColorModeValue, LinkBox, Badge, LinkOverlay, Stack, Checkbox, Code, IconButton, Popover, PopoverArrow, PopoverBody, PopoverCloseButton, PopoverContent, PopoverHeader, PopoverTrigger, Menu, MenuButton, MenuItem, MenuList } from "@chakra-ui/react"
import { useState, useEffect, Fragment } from "react"
import { IoEyeOutline, IoFlashOffOutline } from "react-icons/io5"
import { useDispatch, useSelector } from "react-redux"
import { useSearchParams } from "react-router-dom"
import { Link as RouterLink } from "react-router-dom"
import { Pagination } from "../../components/Pagination"
import { AppDispatch, RootState } from "../../configureStore"
import { Query } from "../../models"
import { getLeaks } from "../../store/slices/leak"

const Leaks = () => {
  const header = ['#', 'Rule', 'Severity', 'File', 'Author', 'Offender', '']
  const dispatch = useDispatch<AppDispatch>()
  const { leaks } = useSelector((state: RootState) => state.leaks)
  const { categories, pagination } = useSelector((state: RootState) => state.global)
  const [searchParams, setSearchParams] = useSearchParams()
  const [loading, setLoading] = useState(false)
  const [query, setQuery] = useState<Query>({
    limit: 10
  })
  useEffect(() => {
    setLoading(true)
    dispatch(getLeaks({
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
      <Table variant="simple" background={useColorModeValue('white', 'gray.800')}>
        <Thead>
          <Tr>
            {header.map((value) => (
              <Th key={value}>{value}</Th>
            ))}
          </Tr>
        </Thead>
        <Tbody>
          {!loading && leaks.length > 0 ? leaks.map(function(leak) {
            return (
              <LinkBox as="tr" key={leak.id}>
                <Td><Checkbox></Checkbox></Td>
                <Td>
                  <LinkOverlay as={RouterLink} to={`/leaks/${leak.id}`}>
                    {leak.rule.display_name}
                  </LinkOverlay>
                  <Stack direction="row">
                    {leak.rule.tags && leak.rule.tags.map(function(tag){
                      return (
                        <Badge variant="outline" colorScheme="brand" key={tag}>
                          {tag}
                        </Badge>
                      )
                    })}
                  </Stack>
                </Td>
                <Td><Badge>{leak.severity}</Badge></Td>
                <Td>
                  <Link href={leak.repository.repository_url + "/commit/" + leak.commit_hash} isExternal>
                    {leak.repository.repository_url}
                  </Link>
                  <Text>{leak.file}</Text>
                  <Text>Line: {leak.line_number}</Text>
                </Td>
                <Td>
                  <Text fontWeight="bold">{leak.author_name}</Text>
                  <Text>{leak.created_at && new Intl.DateTimeFormat("en-GB", {
                    year: "numeric",
                    month: "long",
                    day: "2-digit",
                    hour: "2-digit",
                    minute: "2-digit",
                    second: "2-digit",
                  }).format(Date.parse(leak.created_at))}</Text>
                </Td>
                <Td>
                  <Popover placement="bottom-start">
                    <PopoverTrigger>
                      <IconButton aria-label="Reveal" colorScheme="gray" size="sm" icon={<IoEyeOutline />} />
                    </PopoverTrigger>
                    <PopoverContent>
                      <PopoverArrow />
                      <PopoverCloseButton />
                      <PopoverHeader>Offender</PopoverHeader>
                      <PopoverBody>
                        <Text>{leak.offender}</Text>
                        <Code>{leak.line}</Code>
                      </PopoverBody>
                    </PopoverContent>
                  </Popover>
                </Td>
              </LinkBox>
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

export default Leaks
