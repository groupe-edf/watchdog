import { Icon } from "@chakra-ui/react";
import { Component } from "react";
import { IoPieChartOutline } from "react-icons/io5";
import { ResponsiveContainer, PieChart, Pie, LabelList, Cell, Tooltip } from "recharts";

interface ChartProps {
  data: any
}

class LeaksBySeverity extends Component<ChartProps> {
  constructor(props: ChartProps) {
    super(props)
  }
  render() {
    const { data } = this.props
    const COLORS = ['#f14e32', '#0388a6', '#4e443c']
    return (
      <>
      {data && data.length > 0 ? (
        <ResponsiveContainer width="100%" height="100%">
          <PieChart width={400} height={400}>
            <Pie
              data={data}
              dataKey="y"
              nameKey="x"
              label={false}
              labelLine={false}>
              <LabelList dataKey="x" position="outside" strokeWidth={0} />
              {data && data.map((entry: any, index: number) => (
                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
              ))}
            </Pie>
            <Tooltip />
          </PieChart>
        </ResponsiveContainer>
      ) : ""}
      </>
    )
  }
}

export { LeaksBySeverity }
