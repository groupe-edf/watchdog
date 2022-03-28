import { Button, Drawer, DrawerBody, DrawerContent, DrawerFooter, DrawerHeader, DrawerOverlay, FormControl, FormErrorMessage, FormLabel, Input, Select, Stack, useDisclosure } from "@chakra-ui/react"
import { Fragment } from "react"
import { useForm } from "react-hook-form"
import { IoAddOutline } from "react-icons/io5"
import { useDispatch, useSelector } from "react-redux"
import { AppDispatch, RootState } from "../../configureStore"
import { PolicyCondition } from "../../models"
import { addPolicyCondition } from "../../store/slices/policy"

const AddCondition = (props: { policyId: number }) => {
  const { policyId } = props
  const dispatch = useDispatch<AppDispatch>()
  const { categories } = useSelector((state: RootState) => state.global)
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<PolicyCondition>()
  const { isOpen, onOpen, onClose } = useDisclosure()
  const onSubmit = (condition: PolicyCondition) => {
    condition.policy_id = policyId
    dispatch(addPolicyCondition(condition)).unwrap().then(() => {
      onClose()
    })
  }
  return (
    <Fragment>
      <Button leftIcon={<IoAddOutline/>} colorScheme="brand" onClick={onOpen} size="sm">Add</Button>
      <Drawer isOpen={isOpen} placement="right" onClose={onClose} size="sm">
        <DrawerOverlay />
        <form onSubmit={handleSubmit(onSubmit)}>
        <DrawerContent>
          <DrawerHeader borderBottomWidth="1px">
            Add Condition
          </DrawerHeader>
          <DrawerBody>
            <Stack spacing={6}>
              <FormControl isInvalid={!!errors.type} isRequired>
                <FormLabel htmlFor="type">Type</FormLabel>
                <Select {...register('type', {
                    required: 'Type is required'
                  })}>
                  {categories.filter(category => category.extension === 'condition_type').map((category) => (
                    <option value={category.value}>{category.title}</option>
                  ))}
                </Select>
                <FormErrorMessage>
                  {errors.type && errors.type.message}
                </FormErrorMessage>
              </FormControl>
              <FormControl isInvalid={!!errors.pattern}isRequired>
                <FormLabel htmlFor="pattern">Pattern</FormLabel>
                <Input type="text" {...register('pattern', {
                  required: 'Pattern is required'
                })} />
                <FormErrorMessage>
                  {errors.pattern && errors.pattern.message}
                </FormErrorMessage>
              </FormControl>
            </Stack>
          </DrawerBody>
          <DrawerFooter borderTopWidth="1px">
            <Button variant="outline" mr={3} onClick={onClose}>Cancel</Button>
            <Button type="submit" colorScheme="brand" isLoading={isSubmitting}>Add</Button>
          </DrawerFooter>
        </DrawerContent>
        </form>
      </Drawer>
    </Fragment>
  )
}

export { AddCondition }
