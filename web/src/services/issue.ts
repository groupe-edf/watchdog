import { API_PATH } from "../constants";
import { Issue } from "../store/issues/types";
import { Query, fetchData } from "./commons";



class IssueService {
  async findAll(query: Query) {
    let url = `${API_PATH}/issues?limit=${query?.limit ? query.limit : 10}&offset=${query?.offset ? query.offset : 0}`
    if (query?.conditions.length > 0) {
      for (let condition of query?.conditions) {
        url += `&conditions=${condition.field},${condition.operator},${condition.value}`
      }
    }
    return fetchData<Issue[]>("GET", url);
  }
}

export default new IssueService();
