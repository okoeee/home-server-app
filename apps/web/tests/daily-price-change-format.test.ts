import assert from 'node:assert/strict'
import { test } from 'node:test'
import { formatChangeRate } from '../src/models/daily-price-change.ts'

for (const [rate, expected] of [
  ['10', '10.00'],
  ['2.5', '2.50'],
  ['1.234', '1.23'],
  ['1.236', '1.24'],
  ['1.005', '1.01'],
  ['-1.005', '-1.01'],
  ['-3.456', '-3.46'],
  ['0', '0.00'],
]) {
  test(`騰落率${rate}を小数点以下2桁の${expected}で表示する`, () => {
    // Arrange
    const input = rate

    // Act
    const actual = formatChangeRate(input)

    // Assert
    assert.equal(actual, expected)
  })
}
