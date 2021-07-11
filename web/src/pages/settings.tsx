import { Box, Button, chakra, Checkbox, Flex, GridItem, SimpleGrid, Stack, Text, useColorModeValue } from "@chakra-ui/react"

function Settings() {
  return (
    <SimpleGrid>
    <GridItem mt={[5, null, 0]} colSpan={{ md: 2 }}>
    <Stack
      px={4}
      py={5}
      p={[null, 6]}
      bg={useColorModeValue("white",'gray.700')}
      spacing={6}>
      <chakra.fieldset>
      <Box
        as="legend"
        fontSize="md"
        color={useColorModeValue("gray.900", "gray.50")}>
        Security
      </Box>
        <Stack mt={4} spacing={4}>
          <Flex alignItems="start">
            <Flex alignItems="center" h={5}>
              <Checkbox
                colorScheme="brand"
                id="comments"
                rounded="md"/>
            </Flex>
            <Box ml={3} fontSize="sm">
              <chakra.label
                for="comments"
                fontWeight="md"
                color={useColorModeValue("gray.700", "gray.50")}>
                Reveal Secrets
              </chakra.label>
              <Text color={useColorModeValue("gray.500", "gray.400")}>
                Whether to fully or partially reveal secrets in report and logs.
              </Text>
            </Box>
          </Flex>
        </Stack>
      </chakra.fieldset>
    </Stack>
    <Box
      px={{ base: 4, sm: 6 }}
      py={3}
      bg={useColorModeValue("gray.50", "gray.900")}
      textAlign="right">
      <Button
        type="submit"
        colorScheme="teal"
        fontWeight="md">
        Save
      </Button>
    </Box>
    </GridItem>
    </SimpleGrid>
  )
}

export { Settings }
