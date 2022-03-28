import { useParams } from "react-router-dom"
import { Alert, AlertIcon, Badge, Button, Flex, Grid, GridItem, Heading, HStack, Icon, SkeletonText, Stack, Table, Tbody, Td, Text, Th, Thead, Tr, useToast } from "@chakra-ui/react"
import { useEffect, useState } from "react"
import { useDispatch, useSelector } from "react-redux"
import { AppDispatch, RootState } from "../../configureStore"
import { IntegrationService } from "../../services"
import { IoFlashOffOutline } from "react-icons/io5"
import DeleteIntegration from "./DeleteIntegration"
import { getIntegration } from "../../store/slices/integration"
import Synchronize from "../settings/Synchronize"
import { Card, CardBody } from "@saas-ui/react"

export type IntegrationId = {
  integration_id: string
}
export type IntegrationGroup = {
  id: string
  installed: boolean
  name: string
  path: string
}
const IntegrationView = () => {
  const { integration_id } = useParams<IntegrationId>() as IntegrationId
  const dispatch = useDispatch<AppDispatch>()
  const [ groups, setGroups ] = useState<IntegrationGroup[]>()
  const [loading, setLoading] = useState(false)
  const { integration } = useSelector((state: RootState) => state.integrations)
  const header = ['Gitlab Group', 'Status', 'Actions']
  const toast = useToast()
  useEffect(() => {
    setLoading(true)
    dispatch(getIntegration(integration_id))
    IntegrationService.getGroups(integration_id)
      .then((response) => {
        setGroups(response.data)
      }).catch((error) => {
        toast({
          status: "error",
          title: error.response.data.detail
        })
      }).finally(() => {
        setLoading(false)
      })
  }, [dispatch])
  const installWebhook = async (groupId: string, event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    IntegrationService.installWebhook({
      integration_id: Number(integration_id),
      group_id: groupId
    })
  }
  return (
    <Grid
      gap={4}
      templateColumns='repeat(5, 1fr)'>
      <GridItem colSpan={3}>
        <Card>
          <CardBody>
            <Stack spacing={4}>
            <Heading as='h4' size='md'>Configure your Gitlab integration via group hook</Heading>
            <Alert status='info'>
              <AlertIcon />
              Here is the list of Gitlab groups that Watchdog has access to. Install the ones you want to monitor.
              Note that installing a group will also install all its subgroups
            </Alert>
            <Table variant="simple" size='sm'>
              <Thead>
                <Tr>
                  {header.map((value) => (
                    <Th key={value}>{value}</Th>
                  ))}
                </Tr>
              </Thead>
              <Tbody>
                {groups && groups.length > 0 ? (groups.map(function(group: any) {
                  return (
                    <Tr key={group.id}>
                      <Td>
                        {group.path}
                      </Td>
                      <Td>
                        {group.installed ? <Badge colorScheme='green'>Installed</Badge> : <Badge>Not Installed</Badge>}
                      </Td>
                      <Td>
                        <Button colorScheme="brand" onClick={(event) => installWebhook(group.id, event)} isDisabled={group.installed}>Install</Button>
                      </Td>
                    </Tr>
                  )})) : (
                    (loading ?
                      <Tr key="loading">
                        <Td colSpan={7}>
                          <SkeletonText noOfLines={4} spacing="4" />
                        </Td>
                      </Tr> :
                      <Tr key="empty">
                        <Td colSpan={7} textAlign="center" color="grey" paddingX={3}>
                          <Icon fontSize="64" as={IoFlashOffOutline} />
                          <Text marginTop={4}>No gitlab groups found</Text>
                        </Td>
                      </Tr>
                    )
                  )}
                </Tbody>
            </Table>
            </Stack>
          </CardBody>
        </Card>
      </GridItem>
      <GridItem colSpan={2}>
        <Card>
          <CardBody>
            <HStack>
              <Synchronize integration_id={integration.id}/>
              <DeleteIntegration integrationId={integration.id}/>
            </HStack>
            <Stack spacing={2} marginTop={4}>
              <Flex>
                <Text fontSize='md' fontWeight='bold' marginEnd={4}>
                  Instance Name:
                </Text>
                <Text fontSize='md' fontWeight='400'>
                  {integration.instance_name}
                </Text>
              </Flex>
              <Flex>
                <Text fontSize='md' fontWeight='bold' marginEnd={4}>
                  Instance URL:
                </Text>
                <Text fontSize='md' fontWeight='400'>
                  {integration.instance_url}
                </Text>
              </Flex>
            </Stack>
          </CardBody>
        </Card>
      </GridItem>
    </Grid>
  )
}

export { IntegrationView }
