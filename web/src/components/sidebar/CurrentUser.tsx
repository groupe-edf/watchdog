import { Avatar, Flex, Heading, Text } from "@chakra-ui/react"
import { useSelector } from "react-redux"
import { RootState } from "../../configureStore"

const CurrentUser = (props: any) => {
  const { currentUser } = useSelector((state: RootState) => state.authentication)
  return (
    <Flex align="center">
      <Avatar size="sm" name={currentUser.first_name} />
      <Flex flexDir="column" marginLeft={4}>
        <Heading as="h3" size="sm">{currentUser.last_name}</Heading>
        <Text color="gray">Admin</Text>
      </Flex>
    </Flex>
  )
}

export default CurrentUser
