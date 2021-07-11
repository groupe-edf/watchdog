import { Table, useColorModeValue, TableCaption, Thead, Tr, Th, Tbody, Td, useDisclosure, Box, Button, Stack } from "@chakra-ui/react";
import { FC } from "react";
import { withStatusIndicator } from "../components/withStatusIndicator";
import { API_PATH } from "../constants";
import { useFetch } from "../hooks/useFetch";

export interface Issue {
  id: string,
  author: string,
  commit: Date,
  severity: number,
}

interface IssuesContentProps {
  data: Issue[];
}

export const IssuesContent: FC<IssuesContentProps> = ({ data }) => {
  const header = ['Author', 'Commit', 'Severity'];
  return (
    <>
      <Box paddingBottom={4}>
        <Stack
          justify={'flex-end'}
          direction={'row'}>
          <Button
            colorScheme="teal">
            Purge
          </Button>
        </Stack>
      </Box>
      <Table variant="simple"
        background={useColorModeValue('white', 'gray.800')}>
        <TableCaption>List issues</TableCaption>
        <Thead>
          <Tr>
            {header.map((value) => (
              <Th key={value}>{value}</Th>
            ))}
          </Tr>
        </Thead>
        <Tbody>
          {data && data.map(function(issue){
            return (
              <Tr>
                <Td>{issue.author}</Td>
                <Td>{issue.commit}</Td>
                <Td>{issue.severity}</Td>
              </Tr>
            )
          })}
        </Tbody>
      </Table>
    </>
  )
}
IssuesContent.displayName = 'Issues';
const IssuesWithStatusIndicator = withStatusIndicator(IssuesContent);

function IssuesList() {
  const { isOpen, onOpen, onClose } = useDisclosure()
  const path = `${API_PATH}`;
  const { response, isLoading, error } = useFetch<Issue[]>(`${path}/issues`);
  return (
    <IssuesWithStatusIndicator
      data={response}
      error={error}
      isLoading={isLoading}/>
  )
}

export { IssuesList }
