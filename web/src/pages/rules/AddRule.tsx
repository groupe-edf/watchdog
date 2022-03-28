import { Button, Drawer, DrawerBody, DrawerContent, DrawerFooter, DrawerHeader, DrawerOverlay, FormControl, FormLabel, Input, Select, Stack, Switch, useDisclosure } from "@chakra-ui/react"
import { Fragment, SyntheticEvent, useCallback, useState } from "react"
import { useForm } from "react-hook-form"
import { IoAddOutline } from "react-icons/io5"
import { useDispatch, useSelector } from "react-redux"
import { AppDispatch, RootState } from "../../configureStore"
import { Rule } from "../../models"
import ChakraTagInput from '../../components/ChakraInputTag'
import { addRule, getRules } from "../../store/slices/rule"

const AddRule = () => {
  const { isOpen, onOpen, onClose } = useDisclosure()
  const { categories } = useSelector((state: RootState) => state.global)
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<Rule>()
  const dispatch = useDispatch<AppDispatch>()
  const [tags, setTags] = useState<string[]>([])
  const onSubmit = (rule: Rule) => {
    dispatch(addRule(rule)).unwrap().then(() => {
      onClose()
      dispatch(getRules())
    })
  }
  const handleTagsChange = useCallback((event: SyntheticEvent, tags: string[]) => {
    setTags(tags)
  }, [])
  return (
    <Fragment>
      <Button leftIcon={<IoAddOutline/>} colorScheme="brand" onClick={onOpen}>Add</Button>
      <Drawer isOpen={isOpen} placement="right" onClose={onClose} size="sm">
        <DrawerOverlay />
        <form onSubmit={handleSubmit(onSubmit)}>
        <DrawerContent>
          <DrawerHeader borderBottomWidth="1px">
            Add Rule
          </DrawerHeader>
          <DrawerBody>
            <Stack spacing={4}>
              <FormControl isRequired>
                <FormLabel>Display Name</FormLabel>
                <Input type="text" {...register('display_name', {
                  required: 'Display name is required'
                })}/>
              </FormControl>
              <FormControl isRequired>
                <FormLabel>Pattern</FormLabel>
                <Input type="text" {...register('pattern', {
                  required: 'Pattern is required'
                })}/>
              </FormControl>
              <FormControl>
                <FormLabel>Enabled</FormLabel>
                <Switch defaultChecked={true} colorScheme="brand" {...register('enabled')}/>
              </FormControl>
              <FormControl isRequired>
                <FormLabel>Severity</FormLabel>
                <Select {...register('severity', {
                  required: 'Severity is required'
                })}>
                  {categories.filter(category => category.extension === 'rule_severity').map((category) => (
                    <option value={category.value} key={category.id}>{category.title}</option>
                  ))}
                </Select>
              </FormControl>
              <FormControl>
                <FormLabel>Tags</FormLabel>
                <ChakraTagInput tags={tags} onTagsChange={handleTagsChange}/>
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

export { AddRule }
