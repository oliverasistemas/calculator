export function formatResult(n: number): string {
  if (Number.isInteger(n) && Math.abs(n) < 1e15) {
    return n.toString();
  }
  const s = n.toPrecision(12);
  return parseFloat(s).toString();
}
