import { API_PATH } from "../constants";
import { Rule } from "../store/rules/types";
import { Query, fetchData } from "./commons";

class RuleService {
  async findAll() {
    return fetchData<Rule[]>("GET", `${API_PATH}/rules`);
  }
  async toggle(id: number, isEnabled: boolean) {
    return fetchData<Rule[]>("PUT", `${API_PATH}/rules/${id}/toggle`, {
      enabled: isEnabled
    });
  }
}

export default new RuleService();
