import { Table, Thead, Tr, Th, Tbody, Td, Text, Icon, Box, Stack } from "@chakra-ui/react"
import { Fragment } from "react"
import { IoFlashOffOutline } from "react-icons/io5"
import { AddAPIKey } from "./AddAPIKey"

const APIKeyList = () => {
  const header = ['Name', 'Expires At', 'Revoked', 'Actions']
  return (
    <Fragment>
      <Box paddingBottom={4}>
        <Stack
          justify={'flex-end'}
          direction={'row'}>
          <AddAPIKey/>
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
    </Fragment>
  )
}

export { APIKeyList }
