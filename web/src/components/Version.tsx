import { Component } from "react";
import {
  Box,
  Container,
  HStack,
  Icon,
  Stack,
  Text
} from '@chakra-ui/react';
import { API_PATH } from "../constants";
import { fetchData } from "../services/commons";
import { ApplicationState } from "../store";
import { GlobalActionTypes } from "../store/global/types";
import { connect, ConnectedProps } from "react-redux";
import { RouteComponentProps, withRouter } from "react-router-dom";
import { IoHeartSharp, IoLogoGithub } from "react-icons/io5";

const mapState = (state: ApplicationState) => ({
  state: state.global
})
const mapDispatch = {
  setVersion: (payload: any) => ({ type: GlobalActionTypes.GLOBAL_VERSION, payload }),
}
const connector = connect(mapState, mapDispatch)
type VersionProps = ConnectedProps<typeof connector> & RouteComponentProps

export class Version extends Component<VersionProps, any> {
  constructor(props: VersionProps) {
    super(props)
  }
  componentDidMount() {
    const { state, setVersion } = this.props
    if (Object.keys(state.version).length === 0) {
      fetchData("GET", `${API_PATH}/version`).then(response => {
        setVersion(response.data)
      })
    }
  }
  render() {
    const { state } = this.props
    return (
      <Stack align="center" fontSize="sm" spacing={0}>
        <HStack spacing={1}>
          <Text>Made with</Text>
          <Icon as={IoHeartSharp} color="brand.100" />
          <Text>by</Text>
          <Text fontWeight="bold">Habib MAALEM</Text>
        </HStack>
        <Text>version: {state.version.version}, platform: {state.version.platform}</Text>
      </Stack>
    )
  }
}

export default withRouter(connector(Version));
