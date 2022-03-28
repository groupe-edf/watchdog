import { SimpleGrid, GridItem, Heading, Stack, FormControl, FormLabel, Switch, Select, Input, Text, Textarea, Button, Table, Thead, Tr, Th, Tbody, Td, IconButton, Box, Icon, ButtonGroup, Grid } from "@chakra-ui/react"
import { useEffect } from "react"
import { IoFlashOffOutline, IoTrashOutline } from "react-icons/io5"
import { useDispatch, useSelector } from "react-redux"
import { useNavigate, useParams } from "react-router-dom"
import { Pattern } from "../../components/Pattern"
import { AppDispatch, RootState } from "../../configureStore"
import { deleteCondition, deletePolicy, getPolicy, togglePolicy } from "../../store/slices/policy"
import { AddCondition } from "./AddCondition"
import { Policy, PolicyCondition } from "../../models"
import { useForm } from "react-hook-form"
import { Card, CardBody } from "@saas-ui/react"

export type PolicyId = {
  policy_id: string
}
const EditPolicy = () => {
  const header = ['Type', 'Pattern', '']
  const dispatch = useDispatch<AppDispatch>()
  const navigate = useNavigate()
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm()
  const { categories } = useSelector((state: RootState) => state.global)
  const { policy } = useSelector((state: RootState) => state.policies)
  const { policy_id } = useParams<PolicyId>() as PolicyId
  useEffect(() => {
    dispatch(getPolicy(Number(policy_id)))
  }, [dispatch])
  const handleDeleteCondition = (condition: PolicyCondition) => {
    dispatch(deleteCondition({
      policyId: policy_id,
      condition: condition
    }))
  }
  const handleDeletePolicy = (policy: Policy) => {
    dispatch(deletePolicy(policy))
      .unwrap()
      .then(() => {
        navigate('/policies')
      })
  }
  const onSubmit = (policy: any) => {

  }
  const handleToggle = (policy: Policy) => {
    dispatch(togglePolicy(policy))
  }
  return (
    <Grid
      gap={4}
      templateColumns='repeat(5, 1fr)'>
      <GridItem colSpan={3}>
        <Card>
          <CardBody>
            <Heading fontSize="lg" fontWeight="medium" lineHeight="6" paddingBottom={4}>Conditions</Heading>
            <Table variant="simple">
              <Thead>
                <Tr>
                  {header.map((value) => (
                    <Th key={value}>{value}</Th>
                  ))}
                </Tr>
              </Thead>
              <Tbody>
              {policy.conditions ? policy.conditions.map(function(condition) {
                return (
                <Tr key={condition.id}>
                  <Td>{condition.type}</Td>
                  <Td>{condition.type === 'pattern' ? <Pattern editable={true}>{condition.pattern}</Pattern> : condition.pattern}</Td>
                  <Td textAlign="right"><IconButton aria-label="Delete" size="sm" icon={<IoTrashOutline />} onClick={() => handleDeleteCondition(condition)} /></Td>
                </Tr>
                )
              }) : (
                <Tr key="empty">
                  <Td colSpan={7} textAlign="center" color="grey" paddingX={3}>
                    <Icon fontSize="64" as={IoFlashOffOutline} />
                    <Text marginTop={4}>No conditions found</Text>
                  </Td>
                </Tr>
              )}
              </Tbody>
            </Table>
            <Box paddingTop={3}>
              <AddCondition policyId={Number(policy_id)}/>
            </Box>
          </CardBody>
        </Card>
      </GridItem>
      <GridItem colSpan={2}>
        <Card>
          <CardBody>
          <form onSubmit={handleSubmit(onSubmit)}>
            <Stack paddingBottom={5} spacing={4}>
              <FormControl>
                <FormLabel>Enabled</FormLabel>
                <Switch isChecked={policy.enabled} colorScheme="brand" onChange={() => handleToggle(policy)}/>
              </FormControl>
              <FormControl>
                <FormLabel>Type</FormLabel>
                <Select value={policy.type} isDisabled={true}>
                  {categories.filter(category => category.extension === 'handler_type').map((category) => (
                    <option value={category.value}>{category.title}</option>
                  ))}
                </Select>
              </FormControl>
              <FormControl>
                <FormLabel>Severity</FormLabel>
                <Select value={policy.type}>
                  {categories.filter(category => category.extension === 'issue_severity').map((category) => (
                    <option value={category.value}>{category.title}</option>
                  ))}
                </Select>
              </FormControl>
              <FormControl>
                <FormLabel>Display Name</FormLabel>
                <Input type="text" defaultValue={policy.display_name} {...register('display_name', {
                  required: 'Display name is required'
                })}/>
              </FormControl>
              <FormControl>
                <FormLabel>Description</FormLabel>
                <Textarea {...register('description')} defaultValue={policy.description} />
              </FormControl>
            </Stack>
            <ButtonGroup size='sm' isAttached colorScheme='brand'>
              <Button type="submit">Save</Button>
              <IconButton aria-label='Delete' variant='outline' icon={<IoTrashOutline/>} onClick={() => handleDeletePolicy(policy)} />
            </ButtonGroup>
          </form>
          </CardBody>
        </Card>
      </GridItem>
    </Grid>
  )
}

export { EditPolicy }
