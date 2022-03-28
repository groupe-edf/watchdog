import { Table, useColorModeValue, Thead, Tr, Th, Tbody, Badge, Flex, Icon, Link, SkeletonText, Td, Text, LinkOverlay, LinkBox } from "@chakra-ui/react"
import { useEffect, useState } from "react"
import { IoPerson, IoRocketOutline, IoFlashOffOutline, IoGitCommitOutline, IoTimeOutline } from "react-icons/io5"
import { useDispatch, useSelector } from "react-redux"
import { Link as RouterLink } from "react-router-dom"
import { StatusBadge } from "../../components/StatusBadge"
import { AppDispatch, RootState } from "../../configureStore"
import { Analysis } from "../../models"
import { getAnalyzes } from "../../store/slices/analysis"

const Analyzes = (props: any) => {
  const header = ['State', 'Created By/At', 'Trigger', 'Duration', 'Severity', 'Total Issues']
  const [loading, setLoading] = useState(false)
  const dispatch = useDispatch<AppDispatch>()
  const { analyzes } = useSelector((state: RootState) => state.analyzes)
  useEffect(() => {
    setLoading(true)
    dispatch(getAnalyzes({})).unwrap().then(() => {
      setLoading(false)
    })
  }, [])
  return (
    <Table
      variant="simple"
      background={useColorModeValue('white', 'gray.800')}>
      <Thead>
        <Tr>
          {header.map((value) => (
            <Th key={value}>{value}</Th>
          ))}
        </Tr>
      </Thead>
      <Tbody>
      {analyzes && analyzes.length > 0 ? analyzes.map(function(analysis: Analysis) {
          return (
            <LinkBox as="tr" key={analysis.id}>
              <Td>
                <StatusBadge state={analysis.state} hint={analysis.state_message}/>
              </Td>
              <Td>
                <Flex alignItems="center">
                  <Icon as={IoPerson} marginRight={2}/>
                  <LinkOverlay as={RouterLink} to={`/analyzes/${analysis.id}`}>
                    {analysis.created_by?.first_name} {analysis.created_by?.last_name}
                  </LinkOverlay>
                </Flex>
                <Flex alignItems="center">
                  <Icon as={IoTimeOutline} marginRight={2}/>
                  <Text>
                    {analysis.created_at && new Intl.DateTimeFormat("en-GB", {
                      year: "numeric",
                      month: "long",
                      day: "2-digit",
                      hour: "2-digit",
                      minute: "2-digit",
                      second: "2-digit",
                    }).format(Date.parse(analysis.created_at))}
                  </Text>
                </Flex>
                {analysis.last_commit_hash &&
                <Flex alignItems="center">
                  <Icon as={IoGitCommitOutline} marginRight={2}/>
                  {analysis.last_commit_hash.substring(0, 8)}
                </Flex>
                }
              </Td>
              <Td>
                <Flex alignItems="center">
                  <Icon as={IoRocketOutline} marginRight={2}/>
                  <Badge fontWeight="bold">
                    {analysis.trigger}
                  </Badge>
                </Flex>
              </Td>
              <Td>{analysis.duration && new Date(analysis.duration / 1000 / 1000).toISOString().substr(11, 8)}</Td>
              <Td><StatusBadge state={analysis.severity}/></Td>
              <Td>
                <Link as={RouterLink} to={`/issues?conditions=repository_id,eq,${analysis.repository?.id}`} style={{ textDecoration: 'none' }}>
                  <Text fontWeight="bold" color="brand.100">{analysis.total_issues}</Text>
                </Link>
              </Td>
            </LinkBox>
          )
        }) : (
          (loading ?
            <Tr key="loading">
              <Td colSpan={7}>
                <SkeletonText noOfLines={4} spacing="4" />
              </Td>
            </Tr> :
            <Tr>
              <Td colSpan={header.length} textAlign="center" color="grey" paddingX={4}>
                <Icon fontSize="64" as={IoFlashOffOutline} />
                <Text marginTop={4}>No analyzes found</Text>
              </Td>
            </Tr>
          )
        )}
      </Tbody>
    </Table>
  )
}

export default Analyzes
