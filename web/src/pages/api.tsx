import { Table, Thead, Tr, Th, Box, Button, Stack, Icon, Tbody, Td, Text, Drawer, DrawerContent, DrawerHeader, DrawerOverlay, useDisclosure, DrawerBody, FormLabel, Input, InputGroup, DrawerFooter } from "@chakra-ui/react"
import { Component, useState } from "react"
import { IoAddOutline, IoFlashOffOutline } from "react-icons/io5"
import { connect, ConnectedProps } from "react-redux"
import { RouteComponentProps, withRouter } from "react-router-dom"
import { ApplicationState } from "../store"
import { GlobalActionTypes } from "../store/global/types"
import accessTokenService from "../services/access_token"

export function AddAPIKey() {
  const initialState = {
    api_token: "",
    instance_name: "",
    instance_url: "",
  };
  const { isOpen, onOpen, onClose } = useDisclosure()
  const [values, setValues] = useState(initialState);
  const onChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setValues({ ...values, [event.target.name]: event.target.value })
  }
  const onSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    accessTokenService.save(values).then(data => {
      onClose()
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
            Add API Key
          </DrawerHeader>
          <DrawerBody>
            <Stack spacing="24px">
              <Text
                marginTop={1}
                fontSize="sm"
                color="gray.600">
                You can generate a personal access token for each application you use that needs access to the Watchdog API.
              </Text>
              <Box>
                <FormLabel htmlFor="name">Token name</FormLabel>
                <InputGroup>
                  <Input type="text" name="name" onChange={onChange}/>
                </InputGroup>
              </Box>
              <Box hidden={true}>
                <FormLabel htmlFor="token">Token</FormLabel>
                <InputGroup>
                  <Input type="text" name="token" readOnly={true}/>
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

const mapState = (state: ApplicationState) => ({
  settings: state.global.settings
})
const mapDispatch = {
  getAPIKeys: (payload: any) => ({ type: GlobalActionTypes.GLOBAL_API_KEYS, payload }),
}
const connector = connect(mapState, mapDispatch)
type APIKeyListProps = ConnectedProps<typeof connector> & RouteComponentProps

export class APIKeyList extends Component<APIKeyListProps> {
  constructor(props: APIKeyListProps) {
    super(props)
  }
  componentDidMount() {
    const { settings, getAPIKeys} = this.props
    accessTokenService.findAll().then(response => {
      getAPIKeys(response.data)
    })
  }
  render() {
    const header = ['Name', 'Expires At', 'Revoked', 'Actions']
    const onChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    }
    return (
      <>
      <Box paddingBottom={4}>
        <Stack
          justify={'flex-end'}
          direction={'row'}>
          <AddAPIKey></AddAPIKey>
        </Stack>
      </Box>
      <Table>
        <Thead>
          <Tr>
            {header.map((value) => (
              <Th key={value}>{value}</Th>
            ))}
          </Tr>
        </Thead>
        <Tbody>
          {false ? (
            <Tr>
              <Td></Td>
            </Tr>
          ) : (
            <Tr>
              <Td colSpan={header.length} textAlign="center" color="grey" paddingX={4}>
                <Icon fontSize="64" as={IoFlashOffOutline} />
                <Text marginTop={4}>No API keys found</Text>
              </Td>
            </Tr>
          )}
        </Tbody>
      </Table>
      </>
    )
  }
}

export default withRouter(connector(APIKeyList));
