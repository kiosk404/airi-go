/*
* HTTP client service
*/

export interface ApiResponse<T> {
    code: number;
    message: string;
    data: T;
}

export class ApiError extends Error {
    constructor(public code: number, message: string, public details?: unknown) {
        super(message);
        this.name = 'ApiError';
    }
}

interface RequestOptions extends Omit<RequestInit, 'body'> {
    body?: unknown;
}

class HttpClient {
   private baseUrl: string;

   constructor(baseUrl: string = '') {
      this.baseUrl = baseUrl;
   }

   private async request<T>(
       endpoint: string, options: RequestOptions = {}): Promise<T> {
      const {body, headers, ...rest} = options;

      const config: RequestInit = {
          ...rest,
          headers: {
              'Content-Type': 'application/json',
              ...headers,
          },
      };

      if (body !== undefined) {
          config.body = JSON.stringify(body);
      }

      const url = `${this.baseUrl}${endpoint}`;

      try {
          const response = await fetch(url, config);

          if (!response.ok) {
              const errorText = await response.text();
              throw new ApiError(
                  response.status,
                  errorText || `Unknown Error: ${response.status}`,
                  {url, status: response.status}
              );
          }

          const data = await response.json();

          // 检查业务错误码
          if (data.code !== undefined && data.code !== 0) {
              throw new ApiError(data.code, data.message, data.details);
          }

          return data
      } catch (error) {
          if (error instanceof ApiError) {
              throw error;
          }
          throw new ApiError(
              -1,
              error instanceof Error ? error.message : 'network failed',
              error
          );
      }
   }

   async get<T>(endpoint: string, options?: RequestOptions): Promise<T> {
      return this.request(endpoint, {...options, method: 'GET'});
   }

   async post<T>(endpoint: string, body?: unknown, options?: RequestOptions): Promise<T> {
      return this.request(endpoint, {...options, method: 'POST', body});
   }

   async put<T>(endpoint: string, body?: unknown, options?: RequestOptions): Promise<T> {
      return this.request(endpoint, {...options, method: 'PUT', body});
   }

   async delete<T>(endpoint: string, options?: RequestOptions): Promise<T> {
      return this.request(endpoint, {...options, method: 'DELETE'});
   }
}

export const httpClient = new HttpClient();

export { HttpClient }