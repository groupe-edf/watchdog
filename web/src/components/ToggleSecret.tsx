import { Component } from "react";
import { IoEyeOutline, IoPulseOutline } from "react-icons/io5";
import {
  Box,
  Flex,
  Icon,
} from "@chakra-ui/react";

class ToggleSecret extends Component<any> {
  render() {
    const { secret, occurence } = this.props
    return (
      <Box>
        <Flex alignItems="center">
          <Icon as={IoEyeOutline} marginRight={2}/>
        </Flex>
        <Flex alignItems="center">
          <Icon as={IoPulseOutline} marginRight={2}/>
          {occurence}
        </Flex>
      </Box>
    )
  }
}

export { ToggleSecret };
