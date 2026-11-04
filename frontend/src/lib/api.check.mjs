// Pemeriksaan mandiri untuk extractCode(). Jalankan: node src/lib/api.check.mjs
import assert from 'node:assert/strict'
import { extractCode } from './api.js'

const cases = [
  ['aB3xY9', 'aB3xY9'],
  ['  aB3xY9  ', 'aB3xY9'],
  ['http://localhost:8000/aB3xY9', 'aB3xY9'],
  ['localhost:8000/aB3xY9', 'aB3xY9'],
  ['https://pendek.in/aB3xY9/', 'aB3xY9'],
  ['abc', null],
  ['', null],
  ['http://localhost:8000/', null],
  ['aB3xY9z', null],
  ['aB3-Y9', null],
]

for (const [input, want] of cases) {
  assert.equal(extractCode(input), want, `extractCode(${JSON.stringify(input)})`)
}

console.log(`extractCode: ${cases.length} kasus lolos`)
