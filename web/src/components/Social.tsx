import {
  Button,
  Center,
  Text
} from "@chakra-ui/react";
import { IoLogoGitlab } from "react-icons/io5";

export function Alert(){
  const gitlabAuthorizeUrl = "https://github.com/login/oauth/authorize?client_id=b1f50394c67452fbf5b4&redirect_uri=http://localhost:3001/oauth/redirect"
  return (
    <Button width={'full'} leftIcon={<IoLogoGitlab />} href={gitlabAuthorizeUrl} colorScheme="orange">
      <Center>
        <Text>Continue with Gitlab</Text>
      </Center>
    </Button>
  )
}
