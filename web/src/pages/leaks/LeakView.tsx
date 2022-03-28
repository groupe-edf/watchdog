import { Grid, GridItem, Text, Select, HStack, Link, Code, SimpleGrid } from "@chakra-ui/react"
import { PropertyList, Property, Stepper, StepperStep, Card, CardBody, StepperCompleted, Divider } from "@saas-ui/react"
import { useState, useEffect } from "react"
import { useDispatch, useSelector } from "react-redux"
import { useParams } from "react-router-dom"
import { Commit } from "../../components/Commit"
import { AppDispatch, RootState } from "../../configureStore"
import { getLeakById } from "../../store/slices/leak"

export type LeakId = {
  leak_id: string
}

const LeakView = () => {
  const dispatch = useDispatch<AppDispatch>()
  const { leak_id } = useParams<LeakId>() as LeakId
  const { categories } = useSelector((state: RootState) => state.global)
  const { leak } = useSelector((state: RootState) => state.leaks)
  const [loading, setLoading] = useState(false)
  useEffect(() => {
    setLoading(true)
    dispatch(getLeakById(leak_id)).unwrap().then((response) => {
      setLoading(false)
    })
  }, [])
  return (
    <Grid gap={4} templateColumns='repeat(6, 1fr)'>
      <GridItem colSpan={4}>
        <SimpleGrid gap={4}>
          <Card>
            <CardBody>
              <HStack>
                <Text>Repository</Text>
                <Link href={leak.repository.repository_url + "/commit/" + leak.commit_hash} isExternal>
                  {leak.repository.repository_url}
                </Link>
              </HStack>
            </CardBody>
          </Card>
          <Card>
            <CardBody>
              <Code>{leak.line}</Code>
            </CardBody>
          </Card>
        </SimpleGrid>
      </GridItem>
      <GridItem colSpan={2}>
        <Card>
          <CardBody>
            <PropertyList>
              <Property label="Severity" value={
                <Select value={leak.severity}>
                  {categories.filter(category => category.extension === 'severity').map((category) => (
                    <option key={category.value} value={category.value}>{category.title}</option>
                  ))}
                </Select>
              } />
              <Property label="Committed At" value={leak.created_at && new Intl.DateTimeFormat("en-GB", {
                year: "numeric",
                month: "long",
                day: "2-digit",
                hour: "2-digit",
                minute: "2-digit",
                second: "2-digit",
              }).format(Date.parse(leak.created_at))}/>
              <Property label="Duration" value=""/>
              <Property label="Developer Involved" value={leak.author_name}/>
            </PropertyList>
            <Divider label="How to remediate" marginY={4}/>
            <Stepper orientation="vertical" variant="subtle">
              <StepperStep title="Understand the implications of revoking the secret">
                <Text>A bad situation can be made worse if a secret is revoked without understanding how that secret is currently being used.</Text>
              </StepperStep>
              <StepperStep title="Rotate and revoke the secret">
                <Text>A secret that has been leaked into a git repository should be considered as compromised.</Text>
              </StepperStep>
              <StepperStep title="Remove from git history [Optional]" isActive={true}>
                <Text>If possible, delete the entire repository or rewrite the git history.</Text>
              </StepperStep>
              <StepperStep title="Review access logs" isActive={true}>
                <Text>Check for suspicious activity in the log data of your services impacted by the secret.</Text>
              </StepperStep>
              <StepperCompleted title="You can mark the incident as resolved"/>
            </Stepper>
          </CardBody>
        </Card>
      </GridItem>
    </Grid>
  )
}

export default LeakView
