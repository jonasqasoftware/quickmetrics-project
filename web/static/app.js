const sampleRuns = [
  { id: 'run-101', outcome: 'passed', durationMs: 740, retries: 0 },
  { id: 'run-102', outcome: 'passed', durationMs: 810, retries: 0 },
  { id: 'run-103', outcome: 'failed', durationMs: 1320, retries: 1 },
  { id: 'run-104', outcome: 'passed', durationMs: 920, retries: 0 },
  { id: 'run-105', outcome: 'passed', durationMs: 690, retries: 0 },
]

const form = document.querySelector('#analysis-form')
const runsInput = document.querySelector('#runs')
const results = document.querySelector('#results')
const error = document.querySelector('#error')

function restoreSample() {
  runsInput.value = JSON.stringify(sampleRuns, null, 2)
}

function metric(label, value) {
  const element = document.createElement('div')
  element.className = 'metric'
  const name = document.createElement('span')
  name.textContent = label
  const number = document.createElement('strong')
  number.textContent = value
  element.append(name, number)
  return element
}

function render(data) {
  const grid = document.querySelector('#metric-grid')
  grid.replaceChildren(
    metric('Pass rate', `${data.passRate}%`),
    metric('Failure rate', `${data.failureRate}%`),
    metric('Retry rate', `${data.retryRate}%`),
    metric('Average duration', `${data.averageDurationMs} ms`),
    metric('P95 duration', `${data.p95DurationMs} ms`),
    metric('Total runs', data.totalRuns),
    metric('Passed runs', data.passedRuns),
    metric('Failed runs', data.failedRuns),
  )
  const decision = document.querySelector('#decision')
  decision.textContent = data.decision.status
  decision.className = `decision ${data.decision.status}`
  const reasons = document.querySelector('#reasons')
  reasons.replaceChildren(...data.decision.reasons.map((reason) => {
    const item = document.createElement('li')
    item.textContent = reason
    return item
  }))
  results.hidden = false
  results.focus()
}

form.addEventListener('submit', async (event) => {
  event.preventDefault()
  error.hidden = true
  results.hidden = true
  try {
    const response = await fetch('/api/v1/analyze', {
      method: 'POST',
      headers: { 'content-type': 'application/json' },
      body: JSON.stringify({
        runs: JSON.parse(runsInput.value),
        thresholds: {
          minimumPassRate: Number(document.querySelector('#minimum-pass-rate').value),
          maximumRetryRate: Number(document.querySelector('#maximum-retry-rate').value),
          maximumP95DurationMs: Number(document.querySelector('#maximum-duration').value),
        },
      }),
    })
    const data = await response.json()
    if (!response.ok) throw new Error(data.message || 'Analysis failed')
    render(data)
  } catch (problem) {
    error.textContent = problem instanceof Error ? problem.message : 'Analysis failed'
    error.hidden = false
  }
})

document.querySelector('#load-sample').addEventListener('click', restoreSample)
restoreSample()
