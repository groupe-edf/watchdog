import {
  Badge,
  Button,
  Link,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  TableCaption,
  Text,
  useColorModeValue,
  useDisclosure,
  Box,
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
  IconButton
} from "@chakra-ui/react"
import {
  IoLinkOutline, IoPlayCircleOutline,
} from "react-icons/io5";
import { API_PATH } from '../constants';
import { useFetch } from '../hooks/useFetch';
import { withStatusIndicator } from "../components/withStatusIndicator";
import { FC, useRef, useState } from "react";
import { Route, Switch, useParams, useRouteMatch } from "react-router-dom";
import { RepeatIcon } from "@chakra-ui/icons";

const path = `${API_PATH}`;

export interface Analysis {
  id: string;
  duration?: number;
  finished_at?: string;
  started_at?: string;
  state: string;
  total_issues: number;
  severity: string;
}

export interface Repository {
  id: string;
  last_analysis?: Analysis;
  repository_url: Date;
  issues?: number;
}

interface RepositoriesContentProps {
  data: Repository[];
}

export const RepositoriesContent: FC<RepositoriesContentProps> = ({ data }) => {
  const header = ['URL', 'Last Analysis', 'Status', 'Duration', 'Issues', 'Severity', 'Actions'];
  const match = useRouteMatch();
  return (
    <Table variant="simple"
      background={useColorModeValue('white', 'gray.800')}>
      <TableCaption>List of analyzed repositories</TableCaption>
      <Thead>
        <Tr>
          {header.map((value) => (
            <Th key={value}>{value}</Th>
          ))}
        </Tr>
      </Thead>
      <Tbody>
        {data && data.map(function(repository){
          return (
            <Tr key={repository.id}>
              <Td>
                <Link color="teal.500" to={`${match.url}/${repository.id}`} style={{ textDecoration: 'none' }}>
                  {repository.repository_url}
                </Link>
              </Td>
              <Td>{repository.last_analysis?.started_at && new Intl.DateTimeFormat("en-GB", {
                  year: "numeric",
                  month: "long",
                  day: "2-digit",
                  hour: "2-digit",
                  minute: "2-digit",
                  second: "2-digit",
                }).format(Date.parse(repository.last_analysis?.started_at))}</Td>
              <Td><Badge>{repository.last_analysis?.state}</Badge></Td>
              <Td>{repository.last_analysis?.duration && new Date(repository.last_analysis?.duration / 1000 / 1000).toISOString().substr(11, 8)}</Td>
              <Td>{repository.last_analysis?.total_issues}</Td>
              <Td>{repository.last_analysis?.severity}</Td>
              <Td><IconButton aria-label="Run analysis" colorScheme="teal" icon={<RepeatIcon />} /></Td>
            </Tr>
          )
        })}
      </Tbody>
    </Table>
  )
}
RepositoriesContent.displayName = 'Repositories';
const RepositoriesWithStatusIndicator = withStatusIndicator(RepositoriesContent);

export async function scanCallback(values: any) {
  fetch(`${path}/analyze`, {
    method: 'POST',
    body: JSON.stringify(values),
    headers: { 'Content-Type': 'application/json' }
  })
  .then(response => response.json())
}

export function ScanForm() {
  const { isOpen, onOpen, onClose } = useDisclosure()
  const [isSubmitting, setSubmitting] = useState(false);
  const form = useRef();
  const initialState = {
    repository_url: "",
    from: "",
    since: "",
    until: "",
  };
  const [values, setValues] = useState(initialState);
  const onChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setValues({ ...values, [event.target.name]: event.target.value });
  };
  const onSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    scanCallback(values)
  }

  return (
    <>
      <Box paddingBottom={4}>
        <Stack
          justify={'flex-end'}
          direction={'row'}>
          <Button
            colorScheme="teal"
            onClick={onOpen}>
            New
          </Button>
        </Stack>
      </Box>
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <form onSubmit={onSubmit}>
          <ModalHeader>Scan repository</ModalHeader>
          <ModalBody>
            <Stack spacing={4}>
              <FormControl isRequired>
                <FormLabel htmlFor="repository_url">Repository</FormLabel>
                <InputGroup>
                  <InputLeftElement
                    pointerEvents="none"
                    children={<IoLinkOutline color="gray.300" />}/>
                  <Input type="text" name="repository_url" placeholder="Repository URL" onChange={onChange} />
                </InputGroup>
              </FormControl>
              <FormControl>
                <FormLabel htmlFor="from">Commit</FormLabel>
                <InputGroup>
                  <InputLeftElement
                    pointerEvents="none"
                    children={<IoLinkOutline color="gray.300" />}/>
                  <Input type="text" name="from" placeholder="Commit hash" onChange={onChange} />
                </InputGroup>
              </FormControl>
              <FormControl>
                <FormLabel htmlFor="since">Since</FormLabel>
                <InputGroup>
                  <InputLeftElement
                    pointerEvents="none"
                    children={<IoLinkOutline color="gray.300" />}/>
                  <Input type="date" name="since" placeholder="Since" onChange={onChange} />
                </InputGroup>
              </FormControl>
              <FormControl>
                <FormLabel htmlFor="until">Until</FormLabel>
                <InputGroup>
                  <InputLeftElement
                    pointerEvents="none"
                    children={<IoLinkOutline color="gray.300" />}/>
                  <Input type="date" name="until" placeholder="Until" onChange={onChange} />
                </InputGroup>
              </FormControl>
              <FormControl>
                <FormLabel>Options</FormLabel>
                <Checkbox defaultIsChecked>Enable monitoring</Checkbox>
                <Text color={useColorModeValue("gray.500", "gray.400")}>
                  Frequently scan repositories for secrets
                </Text>
              </FormControl>
            </Stack>
          </ModalBody>
          <ModalFooter>
            <Button
              type="submit"
              isLoading={isSubmitting}
              loadingText="Fetching.."
              colorScheme="teal">
              Analyze
            </Button>
          </ModalFooter>
          </form>
        </ModalContent>
      </Modal>
    </>
  )
}

function RepositoriesList() {
  const { response, isLoading, error } = useFetch<Repository[]>(`${path}/repositories`);
  const match = useRouteMatch();
  return (
    <>
      <ScanForm/>
      <RepositoriesWithStatusIndicator
        data={response}
        error={error}
        isLoading={isLoading}/>
      <Switch>
        <Route path={`${match.path}/:repositoryId`}>
          <Repository />
        </Route>
      </Switch>
    </>
  )
}

export interface RepositoryParams {
  repositoryId: string;
}

function Repository() {
  const { repositoryId } = useParams<RepositoryParams>();
  return (
    <h3>Requested repository Id: {repositoryId}</h3>
  )
}

export { RepositoriesList }
