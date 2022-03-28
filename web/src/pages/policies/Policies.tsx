import { Table, useColorModeValue, Thead, Tr, Th, Tbody, Td, Text, Badge, Flex, Icon, Switch, SkeletonText, LinkBox, LinkOverlay } from "@chakra-ui/react"
import { Fragment, useEffect, useState } from "react"
import {
  Link
} from 'react-router-dom'
import { IoFlashOffOutline, IoMenu } from "react-icons/io5"
import { useDispatch, useSelector } from "react-redux"
import { AppDispatch, RootState } from "../../configureStore"
import { Policy } from "../../models"
import { getPolicies, togglePolicy } from "../../store/slices/policy"
import { AddPolicy } from "./AddPolicy"

const Policies = () => {
  const header = ['Name', 'Type', 'Conditions', 'Severity', 'State']
  const dispatch = useDispatch<AppDispatch>()
  const { policies } = useSelector((state: RootState) => state.policies)
  const [loading, setLoading] = useState(false)
  useEffect(() => {
    setLoading(true)
    dispatch(getPolicies())
      .unwrap()
      .then(() => {
        setLoading(false)
      })
  }, [dispatch])
  const handlePolicyToggle = (policy: Policy) => {
    dispatch(togglePolicy(policy))
  }
  return (
    <Fragment>
      <Flex as="header" align="center" justify="space-between" marginBottom={4} width="full">
        <AddPolicy/>
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
          {policies.length > 0 ? policies.map(function(policy) {
            return (
              <Tr key={policy.id}>
                <Td>
                  <Link to={`/policies/${policy.id}/edit`}>
                    {policy.display_name}
                    {policy.description && (
                      <Text fontSize={'sm'} color="gray">
                        {policy.description}
                      </Text>
                    )}
                  </Link>
                </Td>
                <Td><Badge>{policy.type}</Badge></Td>
                <Td>
                  {policy.conditions &&
                  <Flex alignItems="center">
                    <Icon as={IoMenu} marginRight={2}/>
                    <Text fontWeight="bold">{policy.conditions?.length}</Text>
                  </Flex>
                  }
                </Td>
                <Td><Badge>{policy.severity}</Badge></Td>
                <Td align="right"><Switch isChecked={policy.enabled} colorScheme="brand" onChange={() => handlePolicyToggle(policy)}/></Td>
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

export { Policies }
