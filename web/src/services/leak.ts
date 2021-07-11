import { API_PATH } from "../constants";
import { Leak } from "../store/leaks/types";
import { Query, fetchData } from "./commons";

class LeakService {
  async findAll(query?: Query) {
    let url = `${API_PATH}/leaks?limit=${query?.limit ? query.limit : 10}&offset=${query?.offset ? query.offset : 0}`
    return fetchData<Leak[]>("GET", url);
  }
  async findById(id: string) {
    return fetchData<Leak[]>("GET", `${API_PATH}/leaks/${id}`);
  }
}

export default new LeakService();
