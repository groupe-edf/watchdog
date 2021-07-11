import { useState, useEffect } from 'react';

export type APIResponse<T> = { status: string; data: T };

export interface FetchState<T> {
  response: T;
  error?: Error;
  isLoading: boolean;
}

export const useFetch = <T extends {}>(url: string, options?: RequestInit): FetchState<T> => {
  const [response, setResponse] = useState<T>({ status: 'start fetching' } as any);
  const [error, setError] = useState<Error>();
  const [isLoading, setIsLoading] = useState<boolean>(true);
  useEffect(() => {
    const fetchData = async () => {
      setIsLoading(true);
      try {
        const response = await fetch(url, { cache: 'no-store', credentials: 'same-origin', ...options });
        if (!response.ok) {
          throw new Error(response.statusText);
        }
        const json = (await response.json()) as T;
        setResponse(json);
        setIsLoading(false);
      } catch (error) {
        setError(error);
      }
    };
    fetchData();
  }, [url, options]);
  return { response, error, isLoading };
};
