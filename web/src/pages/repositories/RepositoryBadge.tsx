import { useClipboard, Flex, Image, Input, Button, InputGroup, InputRightElement, Icon } from "@chakra-ui/react"
import { Fragment, useEffect, useState } from "react"
import { IoCopyOutline, IoCheckmarkOutline } from "react-icons/io5"
import { RepositoryService } from "../../services"

const RepositoryBadge = (props: any) => {
  const badgeURL = "http://localhost:3001/api/v1/repositories/" + props.repositoryId + "/badge"
  const [value, setValue] = useState(badgeURL)
  const { hasCopied, onCopy } = useClipboard(value)
  const [badgeData, setBadgeData] = useState<any>()
  useEffect(() => {
    RepositoryService.getBadge(props.repositoryId).then((response) => {
      const data = `data:${response.headers['content-type']};base64,${window.btoa(response.data)}`
      setBadgeData(data)
    })
  }, [])
  return (
    <Fragment>
      <Image src={badgeData} alt="Analyze status"/>
      <Flex marginBlock={2}>
        <InputGroup>
          <Input value={value} isReadOnly paddingRight={2} />
          <InputRightElement>
            <Button onClick={onCopy} width={2} variant="white">
              {hasCopied ? <Icon as={IoCheckmarkOutline} color="green.500"/> : <Icon as={IoCopyOutline}/>}
            </Button>
          </InputRightElement>
        </InputGroup>
      </Flex>
    </Fragment>
  )
}

export default RepositoryBadge
