import { Editable, EditablePreview, EditableInput, Text, IconButton, HStack, Code } from "@chakra-ui/react";
import { Component } from "react";
import { IoFlashOutline, IoTrashOutline } from "react-icons/io5";
import { API_PATH } from "../constants";
import { fetchData } from "../services/commons";

interface PatternProps {
  editable: boolean,
  pattern: string
}

class Pattern extends Component<PatternProps, {
  isSubmitting: boolean
}> {
  constructor(props: PatternProps) {
    super(props);
    this.state = {
      isSubmitting: false
    }
  }
  evaluatePattern = async (event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    event.preventDefault();
    this.setState({isSubmitting: true})
    fetchData("POST", `${API_PATH}/pattern`, {
      pattern: "fsdfsef§§§§"
    }).then(response => {
      this.setState({isSubmitting: false})
    }).catch(response => {
      this.setState({isSubmitting: false})
    })
  }
  render () {
    const {editable, pattern } = this.props
    const { isSubmitting } = this.state
    return (
      <HStack>
      {editable ? (
        <Editable defaultValue="">
          <EditablePreview />
          <EditableInput />
        </Editable>
      ) : (
        <Code>{pattern}</Code>
      )}
      <IconButton
        aria-label="Test"
        isLoading={isSubmitting}
        size="sm"
        variant="outline"
        onClick={this.evaluatePattern}
        icon={<IoFlashOutline />}/>
      </HStack>
    )
  }
}

export { Pattern }
