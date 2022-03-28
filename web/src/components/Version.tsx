import React, { useEffect } from "react"
import { useDispatch, useSelector } from "react-redux"
import {
  HStack,
  Icon,
  Stack,
  Text
} from '@chakra-ui/react'
import { IoHeartSharp } from "react-icons/io5"
import { getVersion } from "../store/slices/global"
import { AppDispatch, RootState } from "../configureStore"

const Version = () => {
  const dispatch = useDispatch<AppDispatch>()
  const { version } = useSelector((state: RootState) => state.global)
  useEffect(() => {
    dispatch(getVersion())
  }, [])
  return (
    <Stack align="center" fontSize="sm" spacing={0}>
      <HStack spacing={1}>
        <Text>Made with</Text>
        <Icon as={IoHeartSharp} color="brand.100" />
        <Text>by</Text>
        <Text fontWeight="bold">Habib MAALEM</Text>
      </HStack>
      {version && <Text>version: {version.version}, platform: {version.platform}, git: {version.git_version}</Text>}
    </Stack>
  )
}

export default React.memo(Version)
