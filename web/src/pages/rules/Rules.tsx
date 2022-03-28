import {
  Table,
  Thead,
  Tr,
  Th,
  Flex,
  Input,
  InputGroup,
  InputLeftElement,
  Tbody,
  Badge,
  Td,
  Switch,
  useColorModeValue
} from "@chakra-ui/react"
import { Fragment, useEffect } from "react"
import { IoSearchOutline } from "react-icons/io5"
import { useDispatch, useSelector } from "react-redux"
import { Tags } from "../../components/Tags"
import { AppDispatch, RootState } from "../../configureStore"
import { getRules } from "../../store/slices/rule"
import { AddRule } from "./AddRule"

const Rules = () => {
  const header = ['Name', 'Severity', 'Tags', 'State']
  const dispatch = useDispatch<AppDispatch>()
  const { rules } = useSelector((state: RootState) => state.rules)
  useEffect(() => {
    dispatch(getRules())
  }, [dispatch])
  const toggleRule = (ruleId: number, enabled: boolean, event: React.ChangeEvent<HTMLInputElement>) => {

  }
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
          <Input name="query" placeholder="Search for rules..."/>
        </InputGroup>
        <AddRule/>
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
          {rules.length > 0 && rules.map(function(rule) {
            return (
              <Tr key={rule.id}>
                <Td>{rule.display_name}</Td>
                <Td><Badge>{rule.severity}</Badge></Td>
                <Td><Tags tags={rule.tags}></Tags></Td>
                <Td><Switch defaultChecked={rule.enabled} colorScheme="brand" onChange={(event) => toggleRule(rule.id, rule.enabled, event)}/></Td>
              </Tr>
            )
          })}
        </Tbody>
      </Table>
    </Fragment>
  )
}
export { Rules }
