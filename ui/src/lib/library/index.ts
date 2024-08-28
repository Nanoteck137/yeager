import { mkdtemp } from "node:fs/promises";
import { join } from "node:path";

export function temp() {
  return mkdtemp(join("/Users/nanoteck137/media/music/temp", "album-"));
}
