import { Stack, Badge } from "@chakra-ui/react"

const Tags = (props: any) => {
  const { tags } = props
  return (
    <>
    {tags && tags.length > 0 &&
    <Stack direction="row">
      {tags.map(function(tag: any){
        return (
          <Badge variant="outline" colorScheme="brand" key={tag}>
            {tag}
          </Badge>
        )
      })}
    </Stack>
    }
    </>
  )
}

export { Tags }
