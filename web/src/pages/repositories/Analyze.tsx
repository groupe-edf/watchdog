import { Button, Collapse, Divider, Drawer, DrawerBody, DrawerContent, DrawerFooter, DrawerHeader, DrawerOverlay, FormControl, FormLabel, IconButton, Input, InputGroup, InputLeftElement, Modal, ModalBody, ModalContent, ModalFooter, ModalHeader, ModalOverlay, Stack, useDisclosure, useToast } from "@chakra-ui/react"
import { Fragment, useState } from "react"
import { useForm } from "react-hook-form"
import { IoLinkOutline, IoPlayOutline, IoPlaySharp } from "react-icons/io5"
import { useDispatch } from "react-redux"
import { AppDispatch } from "../../configureStore"
import { Repository } from "../../models"
import { analyze, analyzeRepository } from "../../store/slices/repository"

const Analyze = (props: any) => {
  const { repository, ...rest } = props
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm()
  const dispatch = useDispatch<AppDispatch>()
  const { isOpen, onOpen, onClose } = useDisclosure()
  const [loading, setLoading] = useState(false)
  const toast = useToast()
  const onSubmit = (values: any) => {
    dispatch(analyze({
      "repository_url": values.repository_url
    })).unwrap().then(() => {
      onClose()
    })
  }
  const handleAnalyzeRepository = (repossitoryId?: string) => {
    if (repossitoryId) {
      dispatch(analyzeRepository({
        "repository_id": repossitoryId
      })).unwrap().then(() => {
        toast({ status: "success", title: "Analysis successfully started" })
      }).catch((error) => {
        toast({ status: "error", title: error })
      })
    }
  }
  return (
    repository ?
    <IconButton
      isLoading={repository.last_analysis && ["in_progress", "pending"].indexOf(repository.last_analysis?.state) > -1}
      aria-label="Run analysis"
      colorScheme="brand"
      icon={<IoPlaySharp />}
      onClick={() => handleAnalyzeRepository(repository?.id)}
      {...rest} /> :
    <Fragment>
      <Button leftIcon={<IoPlayOutline />} isLoading={loading} colorScheme="brand" onClick={onOpen}>
        Analyze
      </Button>
      <Drawer isOpen={isOpen} placement="right" onClose={onClose} size="sm">
        <DrawerOverlay />
        <form onSubmit={handleSubmit(onSubmit)}>
        <DrawerContent>
          <DrawerHeader borderBottomWidth="1px">
            Scan repository
          </DrawerHeader>
          <DrawerBody>
            <Stack spacing={4}>
              <FormControl isRequired>
                <FormLabel htmlFor="repository_url">Repository</FormLabel>
                <InputGroup>
                  <InputLeftElement
                    pointerEvents="none"
                    children={<IoLinkOutline color="gray.300" />}/>
                  <Input type="text" placeholder="Repository URL" {...register('repository_url', {
                    required: 'Repository URL is required'
                  })} />
                </InputGroup>
              </FormControl>
              <FormControl>
                <FormLabel htmlFor="username">Username</FormLabel>
                <InputGroup>
                  <Input type="text" placeholder="Username" {...register('username')} autoComplete="off"/>
                </InputGroup>
              </FormControl>
              <FormControl>
                <FormLabel htmlFor="password">Password</FormLabel>
                <InputGroup>
                  <Input type="password" placeholder="Password" {...register('password')} autoComplete="off"/>
                </InputGroup>
              </FormControl>
              <Divider/>
              <FormControl>
                <FormLabel htmlFor="from">Commit</FormLabel>
                <InputGroup>
                  <InputLeftElement
                    pointerEvents="none"
                    children={<IoLinkOutline/>}/>
                  <Input type="text" placeholder="Commit hash" {...register('from')} />
                </InputGroup>
              </FormControl>
              <FormControl>
                <FormLabel htmlFor="since">Since</FormLabel>
                <InputGroup>
                  <Input type="date" placeholder="Since" {...register('since')} />
                </InputGroup>
              </FormControl>
              <FormControl>
                <FormLabel htmlFor="until">Until</FormLabel>
                <InputGroup>
                  <Input type="date" placeholder="Until" {...register('until')} />
                </InputGroup>
              </FormControl>
            </Stack>
          </DrawerBody>
          <DrawerFooter borderTopWidth="1px">
            <Button variant="outline" mr={3} onClick={onClose}>Cancel</Button>
            <Button type="submit" colorScheme="brand" isLoading={isSubmitting}>Add</Button>
          </DrawerFooter>
        </DrawerContent>
        </form>
      </Drawer>
    </Fragment>
  )
}

export default Analyze
