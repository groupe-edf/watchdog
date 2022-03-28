import { Button } from "@chakra-ui/react"
import { useToast } from "@chakra-ui/toast"
import { useState } from "react"
import { IoSyncOutline } from "react-icons/io5"
import { useDispatch } from "react-redux"
import { useParams } from "react-router-dom"
import { AppDispatch } from "../../configureStore"
import { synchronizeInstance } from "../../store/slices/integration"

export type IntegrationId = {
  integration_id: string
}
const Synchronize = (props: IntegrationId) => {
  const { integration_id } = useParams<IntegrationId>() as IntegrationId
  const dispatch = useDispatch<AppDispatch>()
  const [loading, setLoading] = useState(false)
  const toast = useToast()
  const handleSynchronize = () => {
    setLoading(true)
    dispatch(synchronizeInstance(integration_id)).unwrap().then(() => {
      toast({
        status: "success",
        title: "Integration successfully synced"
      })
    }).catch((error) => {
      toast({
        status: "error",
        title: error.response.data.detail
      })
    }).finally(() => {
      setLoading(false)
    })
  }
  return (
    <Button
      leftIcon={<IoSyncOutline/>}
      onClick={handleSynchronize}
      isLoading={loading}>
      Synchronize
    </Button>
  )
}

export default Synchronize
