import { Component, ReactNode } from "react"
import {
  Container,
  Divider
} from '@chakra-ui/react'
import Sidebar from "./sidebar/Sidebar"
import Version from "./Version"
import { connect } from "react-redux"

type LayoutProps = {
  children?: ReactNode
  title?: string
}

export class Layout extends Component<LayoutProps> {
  constructor(props: LayoutProps) {
    super(props);
  }
  render() {
    const { children } = this.props
    return (
      <Sidebar>
        {children}
        <Divider marginY={4}/>
        <Container
          maxW={'6xl'}
          direction={{ base: 'column', md: 'row' }}
          spacing={4}
          justify={{ md: 'space-between' }}
          align={{ md: 'center' }}>
            <Version />
        </Container>
      </Sidebar>
    )
  }
}

export default connect()(Layout)
