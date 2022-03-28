import {
  Button,
  ButtonGroup,
  Flex,
  HStack,
  Icon,
  IconButton,
  Link,
  SkeletonText,
  Table,
  Tbody,
  Td,
  Text,
  Th,
  Thead,
  Tooltip,
  Tr,
  useToast,
  useColorModeValue
} from "@chakra-ui/react"
import { Fragment, useEffect, useState } from "react"
import { Link as ReactRouterLink, Outlet } from "react-router-dom"
import { IoFlashOffOutline, IoSyncOutline, IoTrashOutline, IoWarningOutline } from "react-icons/io5"
import { useDispatch, useSelector } from "react-redux"
import { AppDispatch, RootState } from "../../configureStore"
import { deleteIntegration, getIntegrations, synchronizeInstance } from "../../store/slices/integration"
import { AddIntegration } from "../integrations/AddIntegration"

const Integrations = () => {
  const header = ['Instance Name', 'Instance URL', 'Created At', 'Last Sync', 'Actions']
  const dispatch = useDispatch<AppDispatch>()
  const [loading, setLoading] = useState(false)
  const { integrations } = useSelector((state: RootState) => state.integrations)
  const toast = useToast()
  useEffect(() => {
    dispatch(getIntegrations())
  }, [dispatch])
  const handleSynchronize = async (integrationId: string, event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    setLoading(true)
    dispatch(synchronizeInstance(integrationId)).unwrap().then(() => {
      setLoading(false)
      toast({
        status: "success",
        title: "Integration successfully synced"
      })
    })
  }
  const handleDelete = async (integrationId: string, event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    dispatch(deleteIntegration(integrationId)).unwrap().then(() => setLoading(false))
  }
  return (
    <Fragment>
      <Flex as="header" align="center" justify="space-between" marginBottom={4} width="full">
        <AddIntegration/>
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
          {integrations.length > 0 ? (integrations.map(function(integration) {
            return (
              <Tr key={integration.id}>
                <Td>
                  <Link as={ReactRouterLink} color="brand.100" to={`/integrations/${integration.id}`} style={{ textDecoration: 'none' }}>
                    {integration.instance_name}
                  </Link>
                </Td>
                <Td>
                  {integration.instance_url}
                </Td>
                <Td>
                  {integration.created_at && new Intl.DateTimeFormat("en-GB", {
                    year: "numeric",
                    month: "long",
                    day: "2-digit",
                    hour: "2-digit",
                    minute: "2-digit",
                    second: "2-digit",
                  }).format(Date.parse(integration.created_at))}
                </Td>
                <Td>
                  <HStack>
                    <Text>
                      {integration.synced_at && new Intl.DateTimeFormat("en-GB", {
                        year: "numeric",
                        month: "long",
                        day: "2-digit",
                        hour: "2-digit",
                        minute: "2-digit",
                        second: "2-digit"
                      }).format(Date.parse(integration.synced_at))}
                    </Text>
                    {integration.syncing_error &&
                    <Tooltip label={integration.syncing_error}>
                      <Icon as={IoWarningOutline} color="brand.100"/>
                    </Tooltip>
                    }
                  </HStack>
                </Td>
                <Td>
                  <IconButton
                    aria-label="Delete"
                    icon={<IoTrashOutline/>}
                    onClick={(event) => handleDelete(integration.id, event)}/>
                </Td>
              </Tr>
            )
          })) : (
            (loading ?
              <Tr key="loading">
                <Td colSpan={7}>
                  <SkeletonText noOfLines={4} spacing="4" />
                </Td>
              </Tr> :
              <Tr key="empty">
                <Td colSpan={7} textAlign="center" color="grey" paddingX={4}>
                  <Icon fontSize="64" as={IoFlashOffOutline} />
                  <Text marginTop={4}>No gitlab integrations found</Text>
                </Td>
              </Tr>
            )
          )}
        </Tbody>
      </Table>
    </Fragment>
  )
}
export { Integrations }
