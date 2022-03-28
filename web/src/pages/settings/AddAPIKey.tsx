import { Alert, AlertIcon, Button, Drawer, DrawerBody, DrawerContent, DrawerFooter, DrawerHeader, DrawerOverlay, FormControl, FormErrorMessage, FormLabel, Input, Stack, useDisclosure } from "@chakra-ui/react"
import { Fragment } from "react"
import { useForm } from "react-hook-form"
import { IoAddOutline } from "react-icons/io5"

const AddAPIKey = () => {
  const { isOpen, onOpen, onClose } = useDisclosure()
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm()
  const onSubmit = (values: any) => {

  }
  return (
    <Fragment>
      <Button leftIcon={<IoAddOutline />} colorScheme="brand" onClick={onOpen}>Add</Button>
      <Drawer isOpen={isOpen} placement="right" onClose={onClose} size="sm">
        <DrawerOverlay />
        <form onSubmit={handleSubmit(onSubmit)}>
        <DrawerContent>
          <DrawerHeader borderBottomWidth="1px">
            Add API Key
          </DrawerHeader>
          <DrawerBody>
            <Stack spacing={4}>
            <Alert status='info'>
              <AlertIcon />
              API Keys can be used to authenticate with watchdog-cli or access the Watchdog API.
            </Alert>
              <FormControl isInvalid={!!errors.name}>
                <FormLabel htmlFor="name">Name</FormLabel>
                <Input type="text" {...register('name', {
                  required: 'Name required'
                })} />
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

export { AddAPIKey }
