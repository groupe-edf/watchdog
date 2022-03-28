import { Table, Thead, Tr, Th, Tbody, Icon, Td, Text, SkeletonText, Flex, Link, Badge } from "@chakra-ui/react"
import { Card, CardBody } from "@saas-ui/react"
import { useEffect, useState } from "react"
import { IoFlashOffOutline, IoGitBranch, IoPerson, IoRocketOutline } from "react-icons/io5"
import { useDispatch, useSelector } from "react-redux"
import { Link as ReactRouterLink } from "react-router-dom"
import { StatusBadge } from "../../components/StatusBadge"
import { AppDispatch, RootState } from "../../configureStore"
import { Analysis } from "../../models"
import { getAnalyzes } from "../../store/slices/analysis"

const LastAnalyzes = (props: any) => {
  const header = ['Started By/At', 'Trigger', 'State']
  const [loading, setLoading] = useState(false)
  const dispatch = useDispatch<AppDispatch>()
  const { analyzes } = useSelector((state: RootState) => state.analyzes)
  useEffect(() => {
    setLoading(true)
    dispatch(getAnalyzes({
      limit: 5,
      sort: [
        {field: "started_at", direction: "desc"}
      ]
    })).unwrap().then(() => {
      setLoading(false)
    })
  }, [])
  return (
    <Card>
      <CardBody>
        <Table>
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
                <Tr key={analysis.id}>
                  <Td>
                    <Flex alignItems="center">
                      <Icon as={IoGitBranch} marginRight={2}/>
                      <Link as={ReactRouterLink} color="brand.100" to={`repositories/${analysis.repository?.id}`} style={{ textDecoration: 'none' }}>
                        {analysis.repository?.repository_url}
                      </Link>
                    </Flex>
                    <Flex alignItems="center">
                      <Icon as={IoPerson} marginRight={2}/>
                      <Text>
                        {analysis.created_by?.first_name} {analysis.created_by?.last_name}
                      </Text>
                    </Flex>
                  </Td>
                  <Td>
                    <Flex alignItems="center">
                      <Icon as={IoRocketOutline} marginRight={2}/>
                      <Badge fontWeight="bold">
                        {analysis.trigger}
                      </Badge>
                    </Flex>
                  </Td>
                  <Td>
                    {analysis.state ? (
                      <StatusBadge state={analysis.state} hint={analysis.state_message} />
                    ) : (
                      ""
                    )}
                  </Td>
                </Tr>
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
      </CardBody>
    </Card>
  )
}

export default LastAnalyzes
