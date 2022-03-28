import {
  Tabs,
  TabList,
  TabPanels,
  Tab,
  TabPanel,
  useColorModeValue
} from "@chakra-ui/react"
import { APIKeyList } from "./APIKeyList"
import { Proxy } from "./Proxy"

const Settings = () => {
  return (
    <Tabs isLazy>
      <TabList>
        <Tab>API Keys</Tab>
        <Tab>Proxy</Tab>
      </TabList>
      <TabPanels background={useColorModeValue('white', 'gray.800')} borderBottomRadius="md">
        <TabPanel><APIKeyList/></TabPanel>
        <TabPanel><Proxy/></TabPanel>
      </TabPanels>
    </Tabs>
  )
}
export { Settings }
