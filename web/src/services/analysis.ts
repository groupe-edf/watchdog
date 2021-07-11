import { API_PATH } from '../constants';
import { Analysis } from '../store/repositories/types';
import { Query, fetchData } from "./commons";

class AnalysisService {
  async deleteById(id: string) {
    return fetchData<Analysis[]>("DELETE", `${API_PATH}/analyzes/${id}`);
  }
  async findAll(query: Query) {
    let url = `${API_PATH}/analyzes?limit=${query?.limit ? query.limit : 10}`
    url += `&offset=${query?.offset ? query.offset : 0}`
    if (query?.sort.length > 0) {
      for (let sort of query?.sort) {
        url += `&sort=${sort.field},${sort.direction}`
      }
    }
    return fetchData<Analysis[]>("GET", url);
  }
  async findAllByRepository(id: string) {
    return fetchData<Analysis[]>("GET", `${API_PATH}/analyzes?conditions=repository_id,eq,${id}`);
  }
  async findById(id: string) {
    return fetchData<Analysis[]>("GET", `${API_PATH}/analyzes/${id}`);
  }
}

export default new AnalysisService();
