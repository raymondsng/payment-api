
const API_BASE = '/api'

async function request(path, options = {}) {
  const response = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers: {
      ...(options.body ? { 'Content-Type': 'application/json' } : {}),
      ...options.headers,
    },
  })

  const text = await response.text()
  let data

  try {
    data = text ? JSON.parse(text) : null
  } catch {
    data = text
  }

  if (!response.ok) {
    const message =
      typeof data === 'object' && data !== null
        ? data.error || data.message || JSON.stringify(data)
        : data || `HTTP ${response.status}`

    throw new Error(`HTTP ${response.status}: ${message}`)
  }

  return data
}

export function createPayment(amountMinor, currency, idempotencyKey) {
  return request('/v1/payments', {
    method: 'POST',
    headers: {
      'Idempotency-Key': idempotencyKey,
    },
    body: JSON.stringify({
      amount_minor: amountMinor,
      currency,
    }),
  })
}

export function getPayment(id) {
  return request(`/v1/payments/${encodeURIComponent(id)}`)
}

export function checkHealth() {
  return request('/healthz')
}
