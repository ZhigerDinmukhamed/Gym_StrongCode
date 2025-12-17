const form = document.getElementById('loginForm');

if (form) {
  form.addEventListener('submit', async (e) => {
    e.preventDefault();

    try {
      const data = await apiFetch('/auth/login', {
        method: 'POST',
        body: JSON.stringify({
          email: email.value,
          password: password.value
        })
      });

      localStorage.setItem('token', data.token);
      window.location.href = 'dashboard.html';

    } catch {
      document.getElementById('error').innerText = 'Login failed';
    }
  });
}

function logout() {
  localStorage.removeItem('token');
  window.location.href = 'index.html';
}
