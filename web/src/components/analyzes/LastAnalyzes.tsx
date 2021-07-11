import {
  Badge,
  Flex,
  Icon,
  Link,
  Table,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
  Text
} from "@chakra-ui/react";
import { Component } from "react";
import { IoPerson, IoRocketOutline, IoGitBranch, IoFlashOffOutline } from "react-icons/io5";
import { Link as ReactRouterLink } from "react-router-dom";
import analysisService from "../../services/analysis"
import { StatusBadge } from "../StatusBadge";
import { Query } from "../../services/commons";

export class LastAnalyzes extends Component<any, {
  analyzes: any,
  query: Query,
  totalItems: number
}> {
  constructor(props: any) {
    super(props)
    this.state = {
      analyzes: [],
      query: {
        conditions: [],
        limit: 5,
        offset: 0,
        sort: [
          {field: "started_at", direction: "desc"}
        ]
      },
      totalItems: 0
    };
  }
  componentDidMount() {
    analysisService.findAll(this.state.query).then(response => this.setState({
      analyzes: response.data
    }))
  }
  render() {
    const header = ['Started By/At', 'Trigger', 'State'];
    return (
      <Table variant="simple" background="white" size="sm">
      <Thead>
        <Tr>
          {header.map((value) => (
            <Th key={value}>{value}</Th>
          ))}
        </Tr>
      </Thead>
      <Tbody>
        {this.state.analyzes.length > 0 ? this.state.analyzes.map(function(analysis: any){
          return (
            <Tr key={analysis.id}>
              <Td>
                <Flex alignItems="center">
                  <Icon as={IoGitBranch} marginRight={2}/>
                  <Link as={ReactRouterLink} color="brand.100" to={`repositories/${analysis.repository.id}`} style={{ textDecoration: 'none' }}>
                    {analysis.repository.repository_url}
                  </Link>
                </Flex>
                <Flex alignItems="center">
                  <Icon as={IoPerson} marginRight={2}/>
                  <Text>
                    {analysis.created_by?.first_name} {analysis.created_by?.last_name}
                  </Text>
                </Flex>
              </Td>
              <Td>
                <Flex alignItems="center">
                  <Icon as={IoRocketOutline} marginRight={2}/>
                  <Badge fontWeight="bold">
                    {analysis.trigger}
                  </Badge>
                </Flex>
              </Td>
              <Td>
                {analysis.state ? (
                  <StatusBadge state={analysis.state} hint={analysis.state_message} />
                ) : (
                  ""
                )}
              </Td>
            </Tr>
          )
        }) : (
          <Tr>
            <Td colSpan={header.length} textAlign="center" color="grey" paddingX={4}>
              <Icon fontSize="64" as={IoFlashOffOutline} />
              <Text marginTop={4}>No analyzes found</Text>
            </Td>
          </Tr>
        )}
      </Tbody>
    </Table>
    )
  }
}
