import { useDispatch } from "react-redux"
import { useForm } from 'react-hook-form'

import { Alert, AlertIcon, Button, Drawer, DrawerBody, DrawerContent, DrawerFooter, DrawerHeader, DrawerOverlay, FormControl, FormErrorMessage, FormLabel, Input, InputGroup, InputRightElement, Stack, Tooltip, useDisclosure } from "@chakra-ui/react"
import { Fragment, useState } from "react"
import { IoAddOutline } from "react-icons/io5"
import { AppDispatch } from "../../configureStore"
import { addIntegration } from "../../store/slices/integration"
import { Integration } from "../../models/integration"

const AddIntegration = () => {
  const { isOpen, onOpen, onClose } = useDisclosure()
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<Integration>()
  const dispatch = useDispatch<AppDispatch>()
  const [show, setShow] = useState(false)
  const handleClick = () => setShow(!show)
  const onSubmit = (values: Integration) => {
    dispatch(addIntegration(values)).unwrap().then(() => {
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
            Add Integration
          </DrawerHeader>
          <DrawerBody>
            <Stack spacing={4}>
              <Alert status='info'>
                <AlertIcon />
                Protect your VCS with both real-time and historical data secrets detection.
              </Alert>
              <FormControl isInvalid={!!errors.instance_url}>
                <FormLabel htmlFor="instance_url">Instance url</FormLabel>
                <Input type="url" {...register('instance_url', {
                  required: 'Instance url is required'
                })} />
                <FormErrorMessage>
                  {errors.instance_url && errors.instance_url.message}
                </FormErrorMessage>
              </FormControl>
              <FormControl>
                <FormLabel htmlFor="instance_name">Name your personal access token</FormLabel>
                <Input type="text" {...register('instance_name', {
                  required: 'Instance name is required'
                })} />
                <FormErrorMessage>
                  {errors.instance_name && errors.instance_name.message}
                </FormErrorMessage>
              </FormControl>
              <FormControl>
                <FormLabel htmlFor="api_token">Personal access token (with api scope)</FormLabel>
                <InputGroup>
                  <Input type={show ? "text" : "password"} {...register('api_token', {
                    required: 'Instance name is required'
                  })}/>
                  <InputRightElement width="4.5rem">
                    <Button h="1.75rem" size="sm" onClick={handleClick}>
                      {show ? "Hide" : "Show"}
                    </Button>
                  </InputRightElement>
                </InputGroup>
                <FormErrorMessage>
                  {errors.api_token && errors.api_token.message}
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

export { AddIntegration }
