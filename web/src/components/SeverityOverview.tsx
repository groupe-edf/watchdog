import { Stat, StatLabel, StatNumber, StatGroup } from "@chakra-ui/react";
import { Component } from "react";

interface AnalyticsProps {
  data: {
    x: string,
    y: string
  }[];
}

class SeverityOverview extends Component<AnalyticsProps> {
  render() {
    const { data } = this.props;
    return (
      <StatGroup
        background="white"
        boxShadow="md"
        padding={4}>
        {data && (data.map(function(item){
          return (
            <Stat>
              <StatLabel>{item.x}</StatLabel>
              <StatNumber>{item.y}</StatNumber>
            </Stat>
          )
        }))}
      </StatGroup>
    )
  }
}

export default SeverityOverview
