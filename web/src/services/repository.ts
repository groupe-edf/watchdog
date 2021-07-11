import { API_PATH } from '../constants';
import { Repository } from '../store/repositories/types';
import { Query, fetchData } from "./commons";

class RepositoryService {
  async analyze(data: any) {
    return fetchData<Repository[]>("POST", `${API_PATH}/analyze`, data);
  }
  async analyzeById(id: string, data: any) {
    return fetchData<Repository[]>("POST", `${API_PATH}/repositories/${id}/analyze`, data);
  }
  async deleteById(id: string) {
    return fetchData<Repository[]>("DELETE", `${API_PATH}/repositories/${id}`);
  }
  async findAll(query?: Query) {
    let url = `${API_PATH}/repositories?limit=${query?.limit ? query.limit : 10}&offset=${query?.offset ? query.offset : 0}&sort=started_at,ASC`
    if (query?.query) {
      url += `&conditions=repository_url,like,${query.query}`
    }
    return fetchData<Repository[]>("GET", url);
  }
  async findById(id: string) {
    return fetchData<Repository[]>("GET", `${API_PATH}/repositories/${id}`);
  }
}

export default new RepositoryService();
