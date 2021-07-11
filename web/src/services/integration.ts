import { API_PATH } from "../constants";
import { Integration } from "../store/integrations/types";
import { Query, fetchData } from "./commons";

class IntegrationService {
  async findAll() {
    return fetchData<Integration[]>("GET", `${API_PATH}/integrations`);
  }
  async findById(id: string) {
    return fetchData<Integration[]>("GET", `${API_PATH}/integrations/${id}`);
  }
  async getRepositories(id: string) {
    return fetchData<Integration[]>("PUT", `${API_PATH}/integrations/${id}/repositories`);
  }
  async save(data: any) {
    return fetchData<Integration[]>("POST", `${API_PATH}/integrations`, data);
  }
  async synchronize(id: string) {
    return fetchData<Integration[]>("GET", `${API_PATH}/integrations/${id}/synchronize`);
  }
}

export default new IntegrationService()
