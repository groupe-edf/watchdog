import { useColorModeValue } from "@chakra-ui/react";

interface TableProps {
  headers: string[],
  items: Object[]
}

export default function Table({ headers, items }: TableProps) {
  return (
    <table></table>
  )
}
