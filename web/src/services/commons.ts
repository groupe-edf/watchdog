export interface Condition {
  field: string
  operator: string
  value?: string
}

export interface Sort {
  field: string
  direction: string
}

export interface Query {
  limit: number
  query?: string
  conditions: Condition[]
  offset: number
  sort: Sort[]
}

export interface Result {
  data: any
  error: any,
  pagination: Range
}

export const fetchData = async <T extends {}>(method: string, url: string, data?: any, options?: RequestInit) => {
  let headers = new Headers({ 'Content-Type': 'application/json' })
  const token = localStorage.getItem("token");
  if (token && token != 'undefined') {
    headers.set("Authorization", ''.concat('Bearer ', token))
  }
  const request = new Request(url, {
    method: method,
    body: JSON.stringify(data),
    headers: headers,
  })
  const response = await fetch(request)
  const result = <Result>{
    data: {},
    error: {},
    pagination: {}
  }
  if (response.status >= 200 && response.status < 300) {
    const conntentRange = response.headers.get("content-range")
    const range = conntentRange && parseContentRange(conntentRange)
    if (range) {
      result.pagination = range
    }
    result["data"] = (await response.json()) as T;
  }
  if (response.status == 401) {
    return Promise.reject({ redirectTo: '/login' })
  }
  if (response.status < 200 || response.status >= 300) {
    return Promise.reject({ error: (await response.json()) })
  }
  return result
};

interface Range {
  unit: string;
  start?: number | null;
  end?: number | null;
  size?: number | null;
}

export function parseContentRange(input: string): Range | null {
  const matches = input.match(/^(\w+) ((\d+)-(\d+)|\*)\/(\d+|\*)$/);
  if (!matches) return null;
  const [, unit, , start, end, size] = matches;
  const range = {
    unit,
    start: start != null ? Number(start) : null,
    end: end != null ? Number(end) : null,
    size: size === "*" ? null : Number(size),
  };
  if (range.start === null && range.end === null && range.size === null)
    return null;
  return range;
}
