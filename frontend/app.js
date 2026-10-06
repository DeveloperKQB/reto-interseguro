// Frontend del reto: consume la API de Go (que a su vez llama a Node).
// Las rutas son relativas (/api/...): nginx las redirige al contenedor de Go,
// así el navegador ve un solo origen y no hace falta CORS.

// El token vive solo en memoria (no en localStorage): si una página sufre
// XSS, no queda un token persistido que robar. Al recargar, se pide login.
let token = null;

const $ = (id) => document.getElementById(id);

function show(el, visible) { el.hidden = !visible; }

function showError(el, message) {
  el.textContent = message;
  show(el, Boolean(message));
}

/** Convierte el texto del textarea en una matriz de números. */
function parseMatrix(text) {
  const rows = text
    .trim()
    .split(/\n+/)
    .map((line) => line.trim())
    .filter(Boolean)
    .map((line) => line.split(/[\s,;]+/).filter(Boolean).map(Number));

  if (rows.length === 0) throw new Error('Escribe al menos una fila.');
  rows.forEach((row, i) => {
    row.forEach((v, j) => {
      if (!Number.isFinite(v)) throw new Error(`Valor no numérico en la fila ${i + 1}, columna ${j + 1}.`);
    });
  });
  return rows;
}

/** Redondea solo para mostrar; el valor exacto queda en el tooltip. */
function fmt(x) {
  return String(Number(x.toFixed(4)));
}

function renderMatrix(el, m) {
  el.replaceChildren();
  el.style.gridTemplateColumns = `repeat(${m[0].length}, auto)`;
  for (const row of m) {
    for (const v of row) {
      const cell = document.createElement('span');
      cell.textContent = fmt(v); // textContent: nunca interpretamos HTML
      cell.title = String(v);
      el.append(cell);
    }
  }
}

async function postJson(url, body, withAuth = false) {
  const headers = { 'Content-Type': 'application/json' };
  if (withAuth) headers.Authorization = `Bearer ${token}`;
  const res = await fetch(url, { method: 'POST', headers, body: JSON.stringify(body) });
  const data = await res.json().catch(() => ({}));
  return { status: res.status, data };
}

function setLoggedIn(loggedIn) {
  show($('login-section'), !loggedIn);
  show($('matrix-section'), loggedIn);
  if (!loggedIn) {
    token = null;
    show($('result-section'), false);
  }
}

$('login-form').addEventListener('submit', async (e) => {
  e.preventDefault();
  const form = new FormData(e.target);
  const button = e.target.querySelector('button');
  button.disabled = true;
  showError($('login-error'), '');
  try {
    const { status, data } = await postJson('/api/auth/login', {
      username: form.get('username'),
      password: form.get('password'),
    });
    if (status !== 200) {
      showError($('login-error'), 'Usuario o contraseña incorrectos.');
      return;
    }
    token = data.token;
    e.target.reset();
    setLoggedIn(true);
  } catch {
    showError($('login-error'), 'No se pudo conectar con el servidor. Revisa que las APIs estén levantadas.');
  } finally {
    button.disabled = false;
  }
});

$('matrix-form').addEventListener('submit', async (e) => {
  e.preventDefault();
  const button = e.target.querySelector('button[type="submit"]');
  showError($('matrix-error'), '');

  let matrix;
  try {
    matrix = parseMatrix($('matrix-input').value);
  } catch (err) {
    showError($('matrix-error'), err.message);
    return;
  }

  button.disabled = true;
  try {
    const { status, data } = await postJson('/api/qr', { matrix }, true);
    if (status === 401) {
      setLoggedIn(false);
      showError($('login-error'), 'Tu sesión expiró. Vuelve a iniciar sesión.');
      return;
    }
    if (status !== 200) {
      showError($('matrix-error'), data.error ?? `Error inesperado (${status}).`);
      return;
    }
    renderMatrix($('out-a'), matrix);
    renderMatrix($('out-q'), data.q);
    renderMatrix($('out-r'), data.r);

    const s = data.stats;
    $('st-max').textContent = fmt(s.max);
    $('st-min').textContent = fmt(s.min);
    $('st-avg').textContent = fmt(s.average);
    $('st-sum').textContent = fmt(s.sum);
    const [q, r] = s.diagonalByMatrix;
    $('st-diag').textContent = `${s.anyDiagonal ? 'Sí' : 'No'} (Q: ${q ? 'sí' : 'no'}, R: ${r ? 'sí' : 'no'})`;
    show($('result-section'), true);
  } catch {
    showError($('matrix-error'), 'No se pudo conectar con el servidor.');
  } finally {
    button.disabled = false;
  }
});

document.querySelectorAll('[data-example]').forEach((btn) => {
  btn.addEventListener('click', () => {
    $('matrix-input').value = btn.dataset.example;
  });
});

$('logout').addEventListener('click', () => setLoggedIn(false));