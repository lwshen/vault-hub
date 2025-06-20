import { PATH } from '@/const/path';
import { AuthApi, Configuration, type ResponseContext } from '@lwshen/vault-hub-ts-fetch-client';
import { navigate } from 'wouter/use-browser-location';

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

        if (!response.ok) {
          let errorMessage = `HTTP ${response.status}: ${response.statusText}`;
          
          try {
            const errorBody = await response.clone().text();
            if (errorBody) {
              try {
                const errorJson = JSON.parse(errorBody);
                errorMessage = errorJson.message || errorJson.error || errorMessage;
              } catch {
                errorMessage = errorBody || errorMessage;
              }
            }
          } catch {
            // Fall back to status text if body parsing fails
          }

          const error = new Error(errorMessage) as ApiError;
          switch (response.status) {
            case 401:
              navigate(PATH.LOGIN);
              break;
            default:
          }
          throw error;
        }
        
        return response;
      },
    },
  ],
});

const authApi = new AuthApi(config);

export { authApi };
