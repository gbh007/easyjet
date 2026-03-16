import prettier from 'eslint-config-prettier';
import vuetify from 'eslint-config-vuetify';

export default (async () => [
  ...(await vuetify({
    ts: true,
  })),
  prettier,
])();
