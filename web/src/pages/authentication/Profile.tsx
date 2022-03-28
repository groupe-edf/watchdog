import { FormControl, FormLabel, Input, Button, VStack, Heading, FormHelperText } from "@chakra-ui/react"
import { Card, CardBody, CardHeader, CardTitle } from "@saas-ui/react"
import { useSelector } from "react-redux"
import { RootState } from "../../configureStore"

const Profile = (props: any) => {
  const { currentUser } = useSelector((state: RootState) => state.authentication)
  const onSubmit = async (event: any) => {

  }
  return (
    <form onSubmit={event => onSubmit(event)}>
      <Card>
        <CardHeader>
          <CardTitle fontSize="xl">Profile</CardTitle>
        </CardHeader>
        <CardBody>
          <VStack spacing={5}>
            <FormControl isRequired>
              <FormLabel htmlFor="first_name">First Name</FormLabel>
              <Input type="first_name" name="first_name" value={currentUser.first_name} placeholder="First Name" />
            </FormControl>
            <FormControl isRequired>
              <FormLabel htmlFor="last_name">Last Name</FormLabel>
              <Input type="last_name" name="last_name" value={currentUser.last_name} placeholder="Last Name" />
            </FormControl>
            <FormControl>
              <FormLabel htmlFor="email">Email</FormLabel>
              <Input type="email" name="email" value={currentUser.email} placeholder="Email" isDisabled={true} />
              <FormHelperText>We currently do not support email address edits</FormHelperText>
            </FormControl>
            <Button type="submit" colorScheme="brand">Update</Button>
          </VStack>
        </CardBody>
      </Card>
    </form>
  )
}

export default Profile
