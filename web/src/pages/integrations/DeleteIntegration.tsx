import { useDisclosure, Button, AlertDialog, AlertDialogBody, AlertDialogContent, AlertDialogFooter, AlertDialogHeader, AlertDialogOverlay } from "@chakra-ui/react"
import React from "react"
import { Fragment } from "react"
import { IoTrashOutline } from "react-icons/io5"
import { useNavigate } from "react-router-dom"
import { IntegrationService } from "../../services"

interface DeleteIntegrationProps {
  integrationId: string
}

const DeleteIntegration = (props: DeleteIntegrationProps) => {
  const { isOpen, onOpen, onClose } = useDisclosure()
  const cancelRef = React.useRef<any>()
  const navigate = useNavigate()
  const handleDelete = () => {
    IntegrationService.delete(props.integrationId).then(() => {
      navigate('/integrations')
    })
  }
  return (
    <Fragment>
      <Button rightIcon={<IoTrashOutline/>} colorScheme="brand" variant="outline" onClick={onOpen}>
        Delete
      </Button>
      <AlertDialog
        isOpen={isOpen}
        leastDestructiveRef={cancelRef}
        onClose={onClose}>
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              Delete Integration
            </AlertDialogHeader>
            <AlertDialogBody>
              Are you sure? You can't undo this action afterwards.
            </AlertDialogBody>
            <AlertDialogFooter>
              <Button onClick={onClose} mr={3}>Cancel</Button>
              <Button colorScheme="brand" onClick={() => handleDelete()}>Confirm</Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>
    </Fragment>
  )
}

export default DeleteIntegration
