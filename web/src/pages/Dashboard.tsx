import { Grid, GridItem } from "@chakra-ui/react"
import LastAnalyzes from "./analyzes/LastAnalyzes"

const Dashboard = () => {
  return (
    <Grid gap={2} templateColumns="repeat(4, 1fr)">
      <GridItem colSpan={2}>
        <LastAnalyzes/>
      </GridItem>
    </Grid>
  )
}
export { Dashboard }
