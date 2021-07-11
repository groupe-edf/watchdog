import { Box, Button, Stack, Table, TableCaption, Tbody, Td, Th, Thead, Tr, useColorModeValue, useDisclosure } from "@chakra-ui/react";
import { FC } from "react";
import { withStatusIndicator } from "../components/withStatusIndicator";
import { API_PATH } from "../constants";
import { useFetch } from "../hooks/useFetch";
import { Issue } from "./issues";

export interface Rule {
  id: string,
  display_name: string,
  enabled: boolean,
  severity: string,
  tags: string[],
}

interface RulesContentProps {
  data: Rule[];
}

export const RulesContent: FC<RulesContentProps> = ({ data }) => {
  const header = ['Name', 'Severity', 'Tags', 'Status'];
  return (
    <>
      <Box paddingBottom={4}>
        <Stack
          justify={'flex-end'}
          direction={'row'}>
          <Button
            colorScheme="teal">
            New
          </Button>
        </Stack>
      </Box>
      <Table variant="simple"
        background={useColorModeValue('white', 'gray.800')}>
        <TableCaption>Rules list</TableCaption>
        <Thead>
          <Tr>
            {header.map((value) => (
              <Th key={value}>{value}</Th>
            ))}
          </Tr>
        </Thead>
        <Tbody>
          {data.map(function(rule){
            return (
              <Tr>
                <Td>{rule.display_name}</Td>
                <Td>{rule.tags}</Td>
                <Td>{rule.enabled}</Td>
              </Tr>
            )
          })}
        </Tbody>
      </Table>
    </>
  )
}
RulesContent.displayName = 'Rules';
const RulesWithStatusIndicator = withStatusIndicator(RulesContent);

function RulesList() {
  const { isOpen, onOpen, onClose } = useDisclosure()
  const path = `${API_PATH}`;
  const { response, isLoading, error } = useFetch<Rule[]>(`${path}/rules`);
  return (
    <RulesWithStatusIndicator
      data={response}
      error={error}
      isLoading={isLoading}/>
  )
}

export { RulesList }
