
import { useState } from 'react'
import { createPayment, getPayment } from './api/payments'
import './App.css'

function App() {
  const [amount, setAmount] = useState('12.50')
  const [currency, setCurrency] = useState('SGD')
  const [paymentId, setPaymentId] = useState('')
  const [result, setResult] = useState(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleCreatePayment(event) {
    event.preventDefault()
    setError('')
    setResult(null)

    const value = Number(amount)

    if (!Number.isFinite(value) || value <= 0) {
      setError('Enter a valid amount greater than zero.')
      return
    }

    // Convert to minor units without sending floating-point currency.
    const amountMinor = Math.round(value * 100)

    if (!Number.isSafeInteger(amountMinor) || amountMinor <= 0) {
      setError('Enter a valid amount.')
      return
    }

    // Generate once per submission; keep the same key for retries
    // of this exact request.
    const idempotencyKey = crypto.randomUUID()

    setLoading(true)

    try {
      const data = await createPayment(
        amountMinor,
        currency,
        idempotencyKey
      )

      setResult(data)
      setPaymentId(data?.id ?? '')
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  async function handleLookupPayment(event) {
    event.preventDefault()
    setError('')
    setResult(null)

    if (!paymentId.trim()) {
      setError('Enter a payment ID.')
      return
    }

    setLoading(true)

    try {
      const data = await getPayment(paymentId.trim())
      setResult(data)
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="page">
      <header className="topbar">
        <div className="brand">
          <span className="brand-icon">P</span>
          <div>
            <h1>PayConsole</h1>
            <p>Payment API dashboard</p>
          </div>
        </div>
        <span className="environment">LOCAL DEVELOPMENT</span>
      </header>

      <section className="intro">
        <p className="eyebrow">PAYMENT MANAGEMENT</p>
        <h2>Payments at a glance.</h2>
        <p>
          Create payment records and retrieve existing payments
          through your Go backend.
        </p>
      </section>

      <section className="cards">
        <article className="panel">
          <div className="panel-heading">
            <span className="panel-icon">＋</span>
            <div>
              <h3>Create payment</h3>
              <p>Submit a new payment request.</p>
            </div>
          </div>

          <form onSubmit={handleCreatePayment}>
            <label htmlFor="amount">Amount</label>
            <div className="amount-input">
              <span>{currency}</span>
              <input
                id="amount"
                type="number"
                min="0.01"
                step="0.01"
                value={amount}
                onChange={(event) => setAmount(event.target.value)}
                required
              />
            </div>

            <label htmlFor="currency">Currency</label>
            <select
              id="currency"
              value={currency}
              onChange={(event) => setCurrency(event.target.value)}
            >
              <option value="SGD">SGD — Singapore Dollar</option>
              <option value="USD">USD — US Dollar</option>
              <option value="EUR">EUR — Euro</option>
            </select>

            <button type="submit" disabled={loading}>
              {loading ? 'Processing...' : 'Create payment'}
            </button>
          </form>

          <p className="hint">
            Demo only. No real money will be transferred.
          </p>
        </article>

        <article className="panel">
          <div className="panel-heading">
            <span className="panel-icon">⌕</span>
            <div>
              <h3>Find a payment</h3>
              <p>Look up a payment by its ID.</p>
            </div>
          </div>

          <form onSubmit={handleLookupPayment}>
            <label htmlFor="paymentId">Payment ID</label>
            <input
              id="paymentId"
              type="text"
              placeholder="Enter payment ID"
              value={paymentId}
              onChange={(event) => setPaymentId(event.target.value)}
              required
            />

            <button
              type="submit"
              className="secondary-button"
              disabled={loading}
            >
              {loading ? 'Loading...' : 'Retrieve payment'}
            </button>
          </form>

          <p className="hint">
            Uses GET /v1/payments/{'{id}'}.
          </p>
        </article>
      </section>

      {error && (
        <section className="message error" role="alert">
          <strong>Request failed</strong>
          <p>{error}</p>
        </section>
      )}

      {result !== null && (
        <section className="panel response-panel">
          <div className="result-heading">
            <div>
              <p className="eyebrow">API RESPONSE</p>
              <h3>Payment details</h3>
            </div>
            <span className="success-label">Response received</span>
          </div>

          <pre>{JSON.stringify(result, null, 2)}</pre>

          {result?.id && (
            <button
              type="button"
              className="secondary-button"
              onClick={() => setPaymentId(result.id)}
            >
              Use this payment ID
            </button>
          )}
        </section>
      )}

      <footer>
        <span>PayConsole</span>
        <span>React + Go + PostgreSQL</span>
      </footer>
    </main>
  )
}

export default App
