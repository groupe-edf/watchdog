import {
  Flex,
  InputGroup,
  InputLeftElement,
  Input,
  Table,
  Tbody,
  Th,
  Thead,
  Tr,
  Switch,
  Td,
  Text,
  Badge,
  Icon,
  LinkBox,
  LinkOverlay,
  useColorModeValue,
  SkeletonText
} from "@chakra-ui/react"
import { Fragment, useEffect, useState } from "react"
import { IoFlashOffOutline, IoLinkOutline, IoMailOutline, IoSearchOutline } from "react-icons/io5"
import { useDispatch, useSelector } from "react-redux"
import { AppDispatch, RootState } from "../../configureStore"
import { User } from "../../models"
import { getUsers } from "../../store/slices/user"

const Users = () => {
  const header = ['Locked', 'Email', 'Full Name', 'Provider', 'Created At', 'Last Login']
  const dispatch = useDispatch<AppDispatch>()
  const [loading, setLoading] = useState(false)
  const { users } = useSelector((state: RootState) => state.users)
  const { currentUser } = useSelector((state: RootState) => state.authentication)
  useEffect(() => {
    dispatch(getUsers())
  }, [dispatch])
  return (
    <Fragment>
      <Flex
        as="header"
        align="center"
        justify="space-between"
        marginBottom={4}
        width="full">
        <InputGroup width="96" display={{ base: "none", md: "flex" }}>
          <InputLeftElement children={<IoSearchOutline/>} />
          <Input name="query" placeholder="Search for users..."/>
        </InputGroup>
      </Flex>
      <Table variant="simple" background={useColorModeValue('white', 'gray.800')}>
        <Thead>
          <Tr>
            {header.map((value) => (
              <Th key={value}>{value}</Th>
            ))}
          </Tr>
        </Thead>
        <Tbody>
          {users.length > 0 ? (users.map(function(user: User) {
            return (
              <Tr key={user.id}>
                <Td><Switch defaultChecked={user.locked} isReadOnly={user.email === currentUser.email} colorScheme="brand"/></Td>
                <Td>
                  <Flex alignItems="center">
                    <Icon as={IoMailOutline} marginRight={2}/>
                    <Text>{user.email}</Text>
                    {user.email === currentUser.email &&
                      <Badge variant="outline" colorScheme="brand" size="xs" marginLeft={2}>Me</Badge>
                    }
                  </Flex>
                  {user.username &&
                    <Flex alignItems="center">
                      <Icon as={IoLinkOutline} marginRight={2}/>
                      <Text>{user.username}</Text>
                    </Flex>
                  }
                </Td>
                <LinkBox as={Td}><LinkOverlay href={`/users/${user.id}/edit`}></LinkOverlay>{user.first_name} {user.last_name}</LinkBox>
                <Td><Badge variant="outline" colorScheme="brand">{user.provider}</Badge></Td>
                <Td>
                  {user.created_at && new Intl.DateTimeFormat("en-GB", {
                    year: "numeric",
                    month: "long",
                    day: "2-digit",
                    hour: "2-digit",
                    minute: "2-digit",
                    second: "2-digit",
                  }).format(Date.parse(user.created_at))}
                </Td>
                <Td>
                  {user.last_login && new Intl.DateTimeFormat("en-GB", {
                    year: "numeric",
                    month: "long",
                    day: "2-digit",
                    hour: "2-digit",
                    minute: "2-digit",
                    second: "2-digit",
                  }).format(Date.parse(user.last_login))}
                </Td>
              </Tr>
            )
          })) : [
            (loading ?
              <Tr key="loading">
                <Td colSpan={7}>
                  <SkeletonText noOfLines={4} spacing="4" />
                </Td>
              </Tr> :
              <Tr key="empty">
                <Td colSpan={7} textAlign="center" color="grey" paddingX={4}>
                  <Icon fontSize="64" as={IoFlashOffOutline} />
                  <Text marginTop={4}>No repositories found</Text>
                </Td>
              </Tr>
            )
          ]}
        </Tbody>
      </Table>
    </Fragment>
  )
}

export { Users }
