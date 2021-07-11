import { API_PATH } from "../constants";
import { Policy } from "../store/policies/types";
import { Query, fetchData } from "./commons";

class PolicyService {
  async findAll() {
    return fetchData<Policy[]>("GET", `${API_PATH}/policies`);
  }
  async findById(id: string) {
    return fetchData<Policy[]>("GET", `${API_PATH}/policies/${id}`);
  }
  async toggle(id: number, isEnabled: boolean) {
    return fetchData<Policy[]>("PUT", `${API_PATH}/policies/${id}/toggle`, {
      enabled: isEnabled
    });
  }
}

export default new PolicyService();
