import http from 'k6/http';
import { sleep } from 'k6';

export let options = { vus: 20, duration: '30s' };

export default function() {
  http.post('https://my-staging.domain/signup', JSON.stringify({
    username: `user${__VU}`,
    password: 'P@ssw0rd!',
    email: `u${__VU}@example.com`
  }), { headers: { 'Content-Type': 'application/json' } });
  sleep(1);
}
