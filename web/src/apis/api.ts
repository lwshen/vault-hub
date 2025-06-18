import { AuthApi, Configuration } from '@lwshen/vault-hub-ts-axios-client';

const config = new Configuration({
  basePath: '',
});

const authApi = new AuthApi(config);

export { authApi };
