import {
  Button,
  Text,
  useColorModeValue,
  useDisclosure,
  Stack,
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Input,
  InputGroup,
  InputLeftElement,
  FormControl,
  FormLabel,
  Checkbox,
  Collapse} from "@chakra-ui/react"
import { Component, useState } from "react"
import { IoLinkOutline, IoPlayOutline } from "react-icons/io5"
import { connect, ConnectedProps } from "react-redux"
import { RouteComponentProps } from "react-router"
import repositoryService from "../../services/repository"
import { ApplicationState } from "../../store"
import { RepositoryActionTypes } from "../../store/repositories/types"

const mapState = (state: ApplicationState) => ({
})
const mapDispatch = {
  appendRepository: (payload: any) => ({ type: RepositoryActionTypes.REPOSITORIES_APPEND, payload }),
}
const connector = connect(mapState, mapDispatch)
type AnalyzeProps = ConnectedProps<typeof connector> & RouteComponentProps

export class Analyze extends Component<AnalyzeProps, any> {
  static INITIAL_STATE = {
    enable_monitoring: false,
    from: "",
    repository_url: "",
    since: "",
    token: "",
    until: "",
    username: "",
  };
  constructor(props: AnalyzeProps) {
    super(props);
    this.state = {
      ...Analyze.INITIAL_STATE,
      isOpen: false,
      isSubmitting: false,
      showAdvanced: false
    }
  }
  static propKey(propertyName: string, value: any): object {
    return { [propertyName]: value };
  }
  setStateWithEvent(event: any, columnType: string): void {
    this.setState(Analyze.propKey(columnType, (event.target as any).value));
  }
  onChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    this.setStateWithEvent(event, event.target.name);
  }
  onSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    this.setState({ isOpen: true })
    event.preventDefault()
    const { appendRepository } = this.props
    repositoryService.analyze(this.state).then(response => {
      appendRepository(response.data)
      this.setState({ isOpen: false })
    })
  }
  render() {
    const { isOpen, isSubmitting, showAdvanced } = this.state
    const handleToggle = () => this.setState({ showAdvanced: !showAdvanced })
    const onClose = () => this.setState({ isOpen: false })
    const onOpen = () => this.setState({ isOpen: true })
    return (
      <>
        <Button
          leftIcon={<IoPlayOutline />}
          isLoading={false}
          loadingText="Analyzing"
          colorScheme="brand"
          onClick={onOpen}
          fontWeight="md">
          Analyze
        </Button>
        <Modal isOpen={isOpen} onClose={onClose}>
          <ModalOverlay />
          <ModalContent>
            <form onSubmit={this.onSubmit}>
            <ModalHeader>Scan repository</ModalHeader>
            <ModalBody>
              <Stack spacing={4}>
                <FormControl isRequired>
                  <FormLabel htmlFor="repository_url">Repository</FormLabel>
                  <InputGroup>
                    <InputLeftElement
                      pointerEvents="none"
                      children={<IoLinkOutline color="gray.300" />}/>
                    <Input type="text" name="repository_url" placeholder="Repository URL" onChange={this.onChange} />
                  </InputGroup>
                </FormControl>
                <FormControl>
                  <FormLabel>Options</FormLabel>
                  <Checkbox defaultIsChecked name="enable_monitoring" onChange={this.onChange}>Enable monitoring</Checkbox>
                  <Text color="gray.500">
                    Frequently scan repositories for secrets
                  </Text>
                </FormControl>
                <Button onClick={handleToggle} marginY={2} colorScheme="brand" variant="outline">Advanced</Button>
                <Collapse in={this.state.showAdvanced} animateOpacity>
                  <FormControl>
                    <FormLabel htmlFor="username">Username</FormLabel>
                    <InputGroup>
                      <Input type="text" name="username" placeholder="Username" onChange={this.onChange} />
                    </InputGroup>
                  </FormControl>
                  <FormControl>
                    <FormLabel htmlFor="token">Token</FormLabel>
                    <InputGroup>
                      <Input type="text" name="token" placeholder="Token" onChange={this.onChange} />
                    </InputGroup>
                  </FormControl>
                  <FormControl>
                    <FormLabel htmlFor="from">Commit</FormLabel>
                    <InputGroup>
                      <InputLeftElement
                        pointerEvents="none"
                        children={<IoLinkOutline color="gray.300" />}/>
                      <Input type="text" name="from" placeholder="Commit hash" onChange={this.onChange} />
                    </InputGroup>
                  </FormControl>
                  <FormControl>
                    <FormLabel htmlFor="since">Since</FormLabel>
                    <InputGroup>
                      <InputLeftElement
                        pointerEvents="none"
                        children={<IoLinkOutline color="gray.300" />}/>
                      <Input type="date" name="since" placeholder="Since" onChange={this.onChange} />
                    </InputGroup>
                  </FormControl>
                  <FormControl>
                    <FormLabel htmlFor="until">Until</FormLabel>
                    <InputGroup>
                      <InputLeftElement
                        pointerEvents="none"
                        children={<IoLinkOutline color="gray.300" />}/>
                      <Input type="date" name="until" placeholder="Until" onChange={this.onChange} />
                    </InputGroup>
                  </FormControl>
                </Collapse>
              </Stack>
            </ModalBody>
            <ModalFooter>
              <Button
                type="submit"
                isLoading={isSubmitting}
                loadingText="Fetching.."
                colorScheme="brand">
                Analyze
              </Button>
            </ModalFooter>
            </form>
          </ModalContent>
        </Modal>
      </>
    )
  }
}

export default connector(Analyze as any)
