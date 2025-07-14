import { PATH } from '@/const/path';
import { AuthApi, Configuration, UserApi, type ResponseContext } from '@lwshen/vault-hub-ts-fetch-client';
import { navigate } from 'wouter/use-browser-location';

interface ApiError extends Error {
  status: number;
  statusText: string;
}

let isNavigatingToLogin = false;
function debounceNavigateToLogin() {
  if (!isNavigatingToLogin) {
    isNavigatingToLogin = true;
    navigate(PATH.LOGIN);
    // Optionally, reset the flag after a short delay if needed:
    setTimeout(() => {
      isNavigatingToLogin = false; 
    }, 1000);
  }
}

const config = new Configuration({
  basePath: '',
  middleware: [
    {
      pre: async (context) => {
        const token = localStorage.getItem('token');
        if (token) {
          context.init.headers = {
            ...context.init.headers,
            Authorization: `Bearer ${token}`,
          };
        }
      },
      post: async (context: ResponseContext) => {
        const { response } = context;

        if (!response.ok) {
          let errorMessage = `HTTP ${response.status}: ${response.statusText}`;
          
          try {
            const errorBody = await response.clone().text();
            if (errorBody) {
              try {
                const errorJson = JSON.parse(errorBody);
                errorMessage = errorJson.error?.message || errorJson.message || errorJson.error || errorMessage;
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
              debounceNavigateToLogin();
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
const userApi = new UserApi(config);

export { authApi, userApi };
