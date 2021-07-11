import moment from "moment";
import { Component } from "react";
import { ResponsiveContainer, Tooltip, Bar, BarChart, CartesianGrid, Label, XAxis, YAxis } from "recharts";

interface ChartProps {
  data: any
}

class LeaksGraph extends Component<ChartProps> {
  constructor(props: ChartProps) {
    super(props)
  }
  render() {
    const { data } = this.props
    return (
      <>
      {data && data.length > 0 ? (
        <ResponsiveContainer width="100%" height="100%">
        <BarChart
          width={400}
          height={400}
          data={data}
          margin={{ top: 20, right: 30, left: 20, bottom: 20 }}>
          <CartesianGrid strokeDasharray="3 3" />
          <Tooltip
            cursor={{ strokeDasharray: '3 3' }}
            formatter={(value: any) =>
              value
            }
            labelFormatter={(label: any) =>
              new Date(label).toLocaleDateString()
            }
          />
          <XAxis dataKey="x" tickFormatter={(date) => moment(date).format('DD-M-YYYY')} name="Date">
            <Label value="Period of time" offset={0} position="bottom" />
          </XAxis>
          <YAxis name="Leaks">
            <Label value="Leak count" angle={-90} position="insideLeft" />
          </YAxis>
          <Bar dataKey="y" fill="#f14e32" />
        </BarChart>
      </ResponsiveContainer>
      ) : ""}
      </>
    )
  }
}

export { LeaksGraph }
