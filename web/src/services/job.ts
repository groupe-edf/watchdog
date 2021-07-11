import { API_PATH } from "../constants";
import { Job } from "../store/jobs/types";
import { Query, fetchData } from "./commons";

class JobService {
  async findAll(query?: Query) {
    let url = `${API_PATH}/jobs?limit=${query?.limit ? query.limit : 10}&offset=${query?.offset ? query.offset : 0}`
    return fetchData<Job[]>("GET", url);
  }
}

export default new JobService();
