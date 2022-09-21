import { Icon } from "@chakra-ui/react";
import { Flex, useColorModeValue, Box, chakra } from "@chakra-ui/react";
import { IoAlertCircleOutline } from "react-icons/io5";

export function Alert(){
  return (
    <Flex
      w="full"
      bg="gray.600"
      p={50}
      alignItems="center"
      justifyContent="center">
      <Flex
        maxW="sm"
        w="full"
        mx="auto"
        bg={useColorModeValue("white", "gray.800")}
        shadow="md"
        rounded="lg"
        overflow="hidden">
        <Flex justifyContent="center" alignItems="center" w={12} bg="red.500">
          <Icon as={IoAlertCircleOutline} color="white" boxSize={6} />
        </Flex>
        <Box mx={-3} py={2} px={4}>
          <Box mx={3}>
            <chakra.span
              color={useColorModeValue("red.500", "red.400")}
              fontWeight="bold">
              Error
            </chakra.span>
            <chakra.p
              color={useColorModeValue("gray.600", "gray.200")}
              fontSize="sm">
              Your email is already used!
            </chakra.p>
          </Box>
        </Box>
      </Flex>
    </Flex>
  );
};

