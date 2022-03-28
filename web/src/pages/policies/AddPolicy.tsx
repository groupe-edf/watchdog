import { useDisclosure, Button, Drawer, DrawerBody, DrawerContent, DrawerHeader, DrawerOverlay, DrawerFooter, FormControl, FormLabel, Input, Select, Stack, Switch, Textarea } from "@chakra-ui/react"
import { Fragment } from "react"
import { useForm } from "react-hook-form"
import { IoAddOutline } from "react-icons/io5"
import { useDispatch, useSelector } from "react-redux"
import { AppDispatch, RootState } from "../../configureStore"
import { Policy } from "../../models"
import { addPolicy } from "../../store/slices/policy"

const AddPolicy = () => {
  const { isOpen, onOpen, onClose } = useDisclosure()
  const { categories } = useSelector((state: RootState) => state.global)
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<Policy>()
  const dispatch = useDispatch<AppDispatch>()
  const onSubmit = (policy: Policy) => {
    dispatch(addPolicy(policy)).unwrap().then(() => {
      onClose()
    })
  }
  return (
    <Fragment>
      <Button leftIcon={<IoAddOutline/>} colorScheme="brand" onClick={onOpen}>Add</Button>
      <Drawer isOpen={isOpen} placement="right" onClose={onClose} size="sm">
        <DrawerOverlay />
        <form onSubmit={handleSubmit(onSubmit)}>
        <DrawerContent>
          <DrawerHeader borderBottomWidth="1px">
            Add Policy
          </DrawerHeader>
          <DrawerBody>
            <Stack paddingBottom={5} spacing={4}>
              <FormControl isRequired>
                <FormLabel>Type</FormLabel>
                <Select {...register('type', {
                  required: 'Type is required'
                })}>
                  {categories.filter(category => category.extension === 'handler_type').map((category) => (
                    <option value={category.value}>{category.title}</option>
                  ))}
                </Select>
              </FormControl>
              <FormControl isRequired>
                <FormLabel>Display Name</FormLabel>
                <Input type="text" {...register('display_name', {
                  required: 'Display name is required'
                })}/>
              </FormControl>
              <FormControl>
                <FormLabel>Description</FormLabel>
                <Textarea {...register('description')} />
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

export { AddPolicy }
