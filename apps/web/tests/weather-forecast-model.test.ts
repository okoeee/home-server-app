import assert from 'node:assert/strict'
import { test } from 'node:test'
import { formatForecastDate } from '../src/models/weather-forecast.ts'

for (const [date, expected] of [
  ['2026-10-31', '10/31（土）'],
  ['2026-11-01', '11/1（日）'],
  ['2026-12-31', '12/31（木）'],
  ['2027-01-01', '1/1（金）'],
  ['2024-02-29', '2/29（木）'],
]) {
  test(`予報日${date}を${expected}と表示する`, () => {
    // Arrange
    const input = date

    // Act
    const actual = formatForecastDate(input)

    // Assert
    assert.equal(actual, expected)
  })
}
