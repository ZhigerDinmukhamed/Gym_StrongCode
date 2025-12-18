const API_BASE = 'http://localhost:8080/api/v1';

function getToken() {
  return localStorage.getItem('token');
}

async function apiFetch(endpoint, options = {}) {
  const res = await fetch(API_BASE + endpoint, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + getToken()
    }
  });

  if (!res.ok) {
    throw new Error('API error');
  }

  return res.json();
}
