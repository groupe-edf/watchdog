import {
  Link} from "@chakra-ui/react"

export function Alert(){
  const gitlabAuthorizeUrl = "https://github.com/login/oauth/authorize?client_id=b1f50394c67452fbf5b4&redirect_uri=http://localhost:3001/oauth/redirect"
  return (
    <Link href={gitlabAuthorizeUrl} isExternal>Continue with Gitlab</Link>
  )
}
