import { AddIcon } from "@chakra-ui/icons"
import { useDisclosure, Button, Drawer, DrawerOverlay, DrawerContent, DrawerCloseButton, DrawerHeader, DrawerBody, Stack, Box, FormLabel, Input, InputGroup, InputLeftAddon, InputRightAddon, Select, Textarea, DrawerFooter } from "@chakra-ui/react"
import { useRef, useState } from "react"
import { IoAddOutline } from "react-icons/io5"
import integrationsService from "../../services/integration"

export function AddIntegration() {
  const initialState = {
    api_token: "",
    instance_name: "",
    instance_url: ""
  }
  const { isOpen, onOpen, onClose } = useDisclosure()
  const [values, setValues] = useState(initialState);
  const onChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setValues({ ...values, [event.target.name]: event.target.value })
  };
  const onSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    integrationsService.save(values).then(data => {
      onClose()
      window.location.reload()
    })
  }
  return (
    <>
      <Button leftIcon={<IoAddOutline />} colorScheme="brand" onClick={onOpen}>Add</Button>
      <Drawer isOpen={isOpen} placement="right" onClose={onClose} size="sm">
        <DrawerOverlay />
        <form onSubmit={onSubmit}>
        <DrawerContent>
          <DrawerHeader borderBottomWidth="1px">
            Add Integration
          </DrawerHeader>
          <DrawerBody>
            <Stack spacing="24px">
              <Box>
                <FormLabel htmlFor="instance_url">Instance url</FormLabel>
                <InputGroup>
                  <Input type="url" name="instance_url" placeholder="https://gitlab.com" onChange={onChange}/>
                </InputGroup>
              </Box>
              <Box>
                <FormLabel htmlFor="instance_name">Name your personal access token</FormLabel>
                <Input type="text" name="instance_name" onChange={onChange}/>
              </Box>
              <Box>
                <FormLabel htmlFor="api_token">Personal access token (with api scope)</FormLabel>
                <InputGroup>
                  <Input type="text" name="api_token" onChange={onChange}/>
                </InputGroup>
              </Box>
            </Stack>
          </DrawerBody>
          <DrawerFooter borderTopWidth="1px">
            <Button variant="outline" mr={3} onClick={onClose}>Cancel</Button>
            <Button type="submit" colorScheme="brand">Add</Button>
          </DrawerFooter>
        </DrawerContent>
        </form>
      </Drawer>
    </>
  )
}
