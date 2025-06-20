import { AuthApi, Configuration, type ResponseContext } from '@lwshen/vault-hub-ts-fetch-client';

interface ApiError extends Error {
  status: number;
  statusText: string;
}

const config = new Configuration({
  basePath: '',
  middleware: [
    {
      post: async (context: ResponseContext) => {
        const { response } = context;
        
        // Handle error responses (4xx, 5xx status codes)
        if (!response.ok) {
          let errorMessage = `HTTP ${response.status}: ${response.statusText}`;
          
          try {
            // Try to parse error response body for more details
            const errorBody = await response.clone().text();
            if (errorBody) {
              try {
                const errorJson = JSON.parse(errorBody);
                errorMessage = errorJson.message || errorJson.error || errorMessage;
              } catch {
                // If not JSON, use the text content
                errorMessage = errorBody || errorMessage;
              }
            }
          } catch {
            // Fall back to status text if body parsing fails
          }
          
          // Throw error with appropriate message
          const error = new Error(errorMessage) as ApiError;
          console.log(error);
          error.status = response.status;
          error.statusText = response.statusText;
          const body = await response.json();
          console.log(body);
          throw error;
        }
        
        return response;
      },
    },
  ],
});

const authApi = new AuthApi(config);

export { authApi };
