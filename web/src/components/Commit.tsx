import { Component } from "react";
import { IoGitCommitOutline, IoLinkOutline } from "react-icons/io5";
import {
  Flex,
  Icon,
  Link,
} from "@chakra-ui/react";

class Commit extends Component<any> {
  render() {
    const { repository, commit } = this.props
    return (
      <>
        <Flex alignItems="center">
          <Icon as={IoLinkOutline} marginRight={2}/>
          <Link href={repository} color="brand.100" isExternal>
            {repository}
          </Link>
        </Flex>
        <Flex alignItems="center">
          <Icon as={IoGitCommitOutline} marginRight={2}/>
          <Link href={repository + "/commit/" + commit.hash} isExternal>
            {commit?.hash.substring(0, 8)}
          </Link>
        </Flex>
      </>
    )
  }
}

export { Commit };
