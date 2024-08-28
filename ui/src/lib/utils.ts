export function ext(s: string) {
  const splits = s.split(".");
  if (splits.length === 1) {
    return "";
  }

  return splits.pop();
}
