import { Component, createRef } from "react";
import {
  AlertDialog,
  AlertDialogBody,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogOverlay,
  Button,
  HStack,
  Icon,
  Stack,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
  Text
} from "@chakra-ui/react";
import { ConnectedProps, connect } from "react-redux";
import { RouteComponentProps } from "react-router-dom";
import repositoryService from "../../services/repository"
import { ApplicationState } from "../../store";
import { RepositoryActionTypes } from "../../store/repositories/types";
import { IoGlobeOutline, IoLockClosedOutline, IoReloadOutline, IoTrashOutline } from "react-icons/io5";
import { AnalyzesList } from ".";
import { StatusBadge } from "../../components/StatusBadge";

const mapState = (state: ApplicationState) => ({
  repository: state.repositories.repository,
})
const mapDispatch = {
  analyze: (payload: any) => ({ type: RepositoryActionTypes.ANALYZE, payload }),
  deleteRepository: (payload: any) => ({ type: RepositoryActionTypes.REPOSITORIES_DELETE_BY_ID, payload }),
  getRepository: (payload: any) => ({ type: RepositoryActionTypes.REPOSITORIES_FIND_BY_ID, payload })
}
interface RepositoryParams {
  repositoryId: string;
}
const connector = connect(mapState, mapDispatch)
type RepositoryProps = ConnectedProps<typeof connector> & RouteComponentProps<RepositoryParams>

class ShowRepository extends Component<RepositoryProps, { isOpen: boolean }> {
  private cancelRef = createRef<any>()
  constructor(props: RepositoryProps) {
    super(props)
    this.deleteRepository = this.deleteRepository.bind(this)
    this.getRepository = this.getRepository.bind(this)
    this.state = {
      isOpen: false
    }
  }
  componentDidMount() {
    const { match } = this.props
    this.getRepository(match.params.repositoryId)
  }
  deleteRepository(id: string, event: React.MouseEvent<HTMLButtonElement, MouseEvent>) {
    const { deleteRepository, history } = this.props
    repositoryService.deleteById(id).then(() => {
      deleteRepository(id)
    })
    history.push('/repositories');
  }
  getRepository(id: string) {
    const { getRepository } = this.props
    repositoryService.findById(id).then(response => {
      getRepository(response.data)
    })
  }
  render() {
    const { analyze, repository } = this.props
    const onClose = () => this.setState({isOpen: false})
    const onClick = async (event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
      event.preventDefault();
      const repository = event.currentTarget.getAttribute("data-repository")
      if (repository != null) {
        repositoryService.analyzeById(repository, {}).then(data => {
          analyze(data)
        })
      }
    }
    return (
      <>
      <Stack spacing={4} direction="row-reverse" align="center">
        <Button
          leftIcon={<IoReloadOutline />}
          isLoading={repository.last_analysis?.state == "STARTED"}
          loadingText="Analyzing"
          colorScheme="brand"
          variant="solid"
          data-repository={repository.id}
          onClick={onClick}>
          Analyze
        </Button>
        <Button rightIcon={<IoTrashOutline />} colorScheme="red" variant="outline"  onClick={() => this.setState({isOpen: true})}>
          Delete
        </Button>
        <AlertDialog
          isOpen={this.state.isOpen}
          leastDestructiveRef={this.cancelRef}
          onClose={onClose}>
          <AlertDialogOverlay>
            <AlertDialogContent>
              <AlertDialogHeader fontSize="lg" fontWeight="bold">
                Delete Repository
              </AlertDialogHeader>
              <AlertDialogBody>
                Are you sure? You can't undo this action afterwards.
              </AlertDialogBody>
              <AlertDialogFooter>
                <Button ref={this.cancelRef} onClick={onClose}>
                  Cancel
                </Button>
                <Button colorScheme="brand" onClick={(event) => this.deleteRepository(repository.id, event)} ml={3}>
                  Confirm
                </Button>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialogOverlay>
        </AlertDialog>
      </Stack>
      <Tabs>
        <TabList>
          <Tab>Overview</Tab>
          <Tab>Analyzes</Tab>
          <Tab isDisabled>Settings</Tab>
        </TabList>
        <TabPanels background="white" borderBottomRadius="md">
          <TabPanel>
            <HStack>
              {repository.visibility === "public" ? (
                <Icon as={IoGlobeOutline} />
              ) : (
                <Icon as={IoLockClosedOutline} />
              )}
              <Text>{repository.repository_url}</Text>
            </HStack>
            <Stack direction="row" marginTop={2}>
              {repository.integration &&
              <StatusBadge state={repository.integration.instance_name}></StatusBadge>
              }
            </Stack>
          </TabPanel>
          <TabPanel>
            {repository.id !== "" &&
            <AnalyzesList repository={repository} />
            }
          </TabPanel>
          <TabPanel>
            <Text>{repository.repository_url}</Text>
          </TabPanel>
        </TabPanels>
      </Tabs>
      </>
    )
  }
}

export default connector(ShowRepository)
