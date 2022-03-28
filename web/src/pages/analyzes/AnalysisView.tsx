import { Stack, Flex, Icon, Text, Badge, Box } from "@chakra-ui/react"
import { IoCalendarOutline, IoGitCommitOutline, IoListOutline, IoPerson, IoRocketOutline, IoTimeOutline } from "react-icons/io5"
import { Analysis } from "../../models"

const AnalysisView = (props: {analysis: Analysis}) => {
  const { analysis } = props
  return (
    <Stack spacing={6}>
      <Flex alignItems="center">
        <Icon as={IoPerson} marginRight={2}/>
        <Text>
          {analysis.created_by?.first_name} {analysis.created_by?.last_name}
        </Text>
      </Flex>
      <Flex alignItems="center">
        <Icon as={IoCalendarOutline} marginRight={2}/>
        <Text>
          {analysis.started_at && new Intl.DateTimeFormat("en-GB", {
            year: "numeric",
            month: "long",
            day: "2-digit",
            hour: "2-digit",
            minute: "2-digit",
            second: "2-digit",
          }).format(Date.parse(analysis.started_at))}
        </Text>
      </Flex>
      <Flex alignItems="center">
        <Icon as={IoGitCommitOutline} marginRight={2}/>
        {analysis.last_commit_hash?.substring(0, 8)}
      </Flex>
      <Flex alignItems="center">
        <Icon as={IoRocketOutline} marginRight={2}/>
        <Badge fontWeight="bold">{analysis.trigger}</Badge>
      </Flex>
      <Flex alignItems="center">
        <Icon as={IoTimeOutline} marginRight={2}/>
        <Text>{analysis.duration && new Date(analysis.duration / 1000 / 1000).toISOString().substr(11, 8)}</Text>
      </Flex>
      <Flex alignItems="center">
        <Icon as={IoListOutline} marginRight={2}/>
        <Text fontWeight="bold" color="brand.100">{analysis.total_issues}</Text>
      </Flex>
    </Stack>
  )
}

export default AnalysisView
