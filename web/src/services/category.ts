import { API_PATH } from "../constants";
import { Category } from "../store/categories/types";
import { fetchData } from "./commons";


class CategoryService {
  async findAll() {
    return fetchData<Category[]>("GET", `${API_PATH}/categories`)
  }
}

export default new CategoryService()
