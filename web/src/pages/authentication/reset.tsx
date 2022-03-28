import {
  Box,
  Heading,
  Link,
  Stack,
  Text
} from "@chakra-ui/react";
import { Component } from "react";
import { Link as RouterLink } from "react-router-dom";
import { connect } from "react-redux";
import { Card } from "../../components/Card";

export class Reset extends Component<any> {
  render() {
    return (
      <Box
        background="gray.50"
        minH="100vh"
        paddingY="12"
        paddingX={{ base: '4', lg: '8' }}>
        <Box maxWidth="md" marginX="auto">
          <Heading textAlign="center" size="xl" fontWeight="bold">
            Forgot your password ?
          </Heading>
          <Text marginTop="4" marginBottom="8" align="center" maxW="md" fontWeight="medium">
            <Text as="span">Did you remember ?</Text>
            <Link as={RouterLink} to="/login" color="brand.100" marginStart="1" display={{ base: 'block', sm: 'inline' }}>Login</Link>
          </Text>
          <Card textAlign="center">
            <Stack spacing={6}>
              <Text>Keep calm and try to remember your password.</Text>
            </Stack>
          </Card>
        </Box>
      </Box>
    )
  }
}

export default connect()(Reset)
