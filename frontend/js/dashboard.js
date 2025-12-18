async function loadClasses() {
    const classes = await apiFetch('/classes');
    const content = document.getElementById('content');
  
    content.innerHTML = classes.map(c => `
      <div class="card">
        <h3>${c.name}</h3>
        <p>${c.description}</p>
        <button onclick="bookClass(${c.id})">Book</button>
      </div>
    `).join('');
  }
  
  async function bookClass(classId) {
    await apiFetch('/bookings', {
      method: 'POST',
      body: JSON.stringify({ class_id: classId })
    });
    alert('Booked!');
  }
  
  async function loadGyms() {
    const gyms = await apiFetch('/gyms');
    content.innerHTML = gyms.map(g => `
      <div class="card">
        <h3>${g.name}</h3>
        <p>${g.address}</p>
      </div>
    `).join('');
  }
  
  async function loadBookings() {
    const bookings = await apiFetch('/bookings');
    content.innerHTML = bookings.map(b => `
      <div class="card">
        Booking #${b.id} — ${b.status}
      </div>
    `).join('');
  }
  