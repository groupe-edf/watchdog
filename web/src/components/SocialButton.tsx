import { HStack, Button, Link, VisuallyHidden } from "@chakra-ui/react"
import { IoLogoGithub } from "react-icons/io5"

const SocialButton = () => {
  return (
    <HStack>
      <Button
        colorScheme="gray"
        rounded={'full'}
        display={'inline-flex'}
        alignItems={'center'}
        justifyContent={'center'}
        size="sm"
        as={'a'}
        href={'https://github.com/groupe-edf/watchdog'}>
        <VisuallyHidden>Source</VisuallyHidden>
        <IoLogoGithub/>
      </Button>
    </HStack>
  )
}

export default SocialButton
