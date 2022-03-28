import { Alert, AlertDescription, AlertIcon, Button, Code, Editable, EditablePreview, EditableTextarea, FormControl, FormLabel, HStack, IconButton, Input, Modal, ModalBody, ModalContent, ModalFooter, ModalHeader, ModalOverlay, useDisclosure } from "@chakra-ui/react"
import { useState } from "react"
import { useForm } from "react-hook-form"
import { IoPlayOutline } from "react-icons/io5"
import { GlobalService } from "../services"

const Pattern = (props: any) => {
  const { children, editable } = props
  const { isOpen, onOpen, onClose } = useDisclosure()
  const [loading, setLoading] = useState(false)
  const [matches, setMatches] = useState()
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm()
  const onSubmit = (values: any) => {
    setLoading(true)
    GlobalService.evaluatePattern(values)
      .then((response) => {
        setLoading(false)
        setMatches(response.data)
      })
  }
  return (
    children &&
      <HStack>
      {editable ? (
        <Editable defaultValue={children}>
        <EditablePreview />
        <EditableTextarea />
      </Editable>
      ) : (
        <Code>{children}</Code>
      )}
      <IconButton
        aria-label="Test"
        size="sm"
        variant="outline"
        onClick={onOpen}
        title="Test"
        icon={<IoPlayOutline />}/>
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <form  onSubmit={handleSubmit(onSubmit)}>
            <ModalHeader>Test a content against your pattern</ModalHeader>
            <ModalBody>
              <Code>{children}</Code>
              <FormControl isRequired>
                <FormLabel htmlFor="payload">Payload</FormLabel>
                <Input type="text" {...register('payload', {
                  required: 'Payload is required'
                })} />
              </FormControl>
              <Input type="hidden" value={children} {...register('pattern')} />
            </ModalBody>
            <ModalFooter>
              <Button
                type="submit"
                isLoading={loading}
                loadingText="Evaluating.."
                colorScheme="brand">
                Evaluate
              </Button>
            </ModalFooter>
          </form>
        </ModalContent>
      </Modal>
    </HStack>
  )
}

export { Pattern }
