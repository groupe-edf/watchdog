import { API_PATH } from "../constants"
import { AccessToken } from "../store/global/types"
import { fetchData } from "./commons"

class AccessTokensService {
  async findAll() {
    return fetchData<AccessToken[]>("GET", `${API_PATH}/access_tokens`)
  }
  async save(data: any) {
    return fetchData<AccessToken[]>("POST", `${API_PATH}/access_tokens`, data);
  }
}

export default new AccessTokensService()
