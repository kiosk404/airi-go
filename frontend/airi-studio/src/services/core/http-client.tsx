/*
* HTTP client service
*/

import JSONBIG from 'json-bigint';

const JSONBigString = JSONBIG({
    storeAsString: true,
});

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

          const responseText = await response.text();

          if (!response.ok) {
              let errorMessage = `Unknown Error: ${response.status}`;
              try {
                  const errorData = await JSONBigString.parse(responseText);
                  errorMessage = errorData.Msg || errorData.msg || errorData.message || errorMessage;
              } catch {
                  // 解析失败，使用原始JSON
                  errorMessage = responseText || errorMessage;
              }
              throw new ApiError(
                  response.status,
                  errorMessage,
                  {url, status: response.status}
              );
          }

          const data = await JSONBigString.parse(responseText);

          // 检查业务错误码（支持 Code 和 code 两种格式）
          const errorCode = data.Code !== undefined ? data.Code : data.code;
          const errorMessage = data.Msg !== undefined ? data.Msg : data.message;

          if (errorCode !== undefined && errorCode !== 0) {
              throw new ApiError(errorCode, errorMessage, data.details);
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