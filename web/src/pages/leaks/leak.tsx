import {
  Box,
  Divider,
  Flex,
  HStack,
  Link,
  Text
} from "@chakra-ui/layout";
import { Component } from "react";
import { connect, ConnectedProps } from "react-redux";
import { RouteComponentProps } from "react-router-dom";
import { Card } from "../../components/Card";
import { ApplicationState } from "../../store";
import { LeakActionTypes } from "../../store/leaks/types";
import leakService from "../../services/leak"
import { Badge, Code, Heading, Stack } from "@chakra-ui/react";
import { ExternalLinkIcon } from "@chakra-ui/icons";

const mapState = (state: ApplicationState) => ({
  leak: state.leaks.leak
})
const mapDispatch = {
  getLeak: (payload: any) => ({ type: LeakActionTypes.LEAKS_FIND_BY_ID, payload }),
}
interface LeakParams {
  leakId: string;
}
const connector = connect(mapState, mapDispatch)
type ShowLeakProps = ConnectedProps<typeof connector> & RouteComponentProps<LeakParams>

export class ShowLeak extends Component<ShowLeakProps> {
  constructor(props: ShowLeakProps) {
    super(props)
  }
  componentDidMount() {
    const { match } = this.props
    this.getLeak(match.params.leakId)
  }
  getLeak(id: string) {
    const { getLeak } = this.props
    leakService.findById(id).then(response => {
      getLeak(response.data)
    })
  }
  render() {
    const { leak } = this.props
    return (
      <Flex>
        <Card flex="2" paddingX={4} paddingY={4}>
          <Text fontSize="3xl" fontWeight="bold" lineHeight="tight">
            {leak.rule.display_name}
          </Text>
          <Stack direction="row">
            {leak.rule.tags && leak.rule.tags.map(function(tag){
              return (
                <Badge variant="outline" colorScheme="brand" key={tag}>
                  {tag}
                </Badge>
              )
            })}
          </Stack>
          <Divider marginY={4}/>
          <Text>{leak.commit_hash} {leak.file}</Text>
          <Text>{leak.offender}</Text>
          <Code>{leak.line}</Code>
        </Card>
        <Box flex="1" padding={4}>
          <Text>A secret has been exposed in your git history</Text>
          <HStack>
            <Text>Repository</Text>
            <Link href={leak.repository.repository_url + "/commit/" + leak.commit_hash} color="brand.100" isExternal>
              {leak.repository.repository_url} <ExternalLinkIcon mx="2px" />
            </Link>
          </HStack>
          <HStack>
            <Text>Severity</Text>
            <Badge variant="outline" colorScheme="brand">{leak.severity}</Badge>
          </HStack>
          <HStack>
            <Text>Developer Involved</Text>
            <Text fontWeight="bold">{leak.author}</Text>
          </HStack>
          <HStack>
            <Text>Pushed at</Text>
            <Text>{leak.created_at && new Intl.DateTimeFormat("en-GB", {
              year: "numeric",
              month: "long",
              day: "2-digit",
              hour: "2-digit",
              minute: "2-digit",
              second: "2-digit",
            }).format(Date.parse(leak.created_at))}</Text>
          </HStack>
        </Box>
      </Flex>
    )
  }
}

export default connector(ShowLeak)
