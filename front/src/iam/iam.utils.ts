import decode from 'jwt-decode';

const claim = 'https://hasura.io/jwt/claims';

export function getUserId(token: string) {
  const decoded = decode<any>(token);

  return decoded[claim]?.['x-hasura-user-id'] ?? '';
}

export function isAdmin(token: string) {
  const decoded = decode<any>(token);

  return decoded[claim]?.['x-hasura-default-role'] === 'admin';
}
