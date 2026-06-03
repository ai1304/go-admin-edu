const DEFAULT_COVERS = Array.from({ length: 6 }, (_, index) => `/images/default-covers/cover-${index + 1}.png`);

function hashSeed(seed) {
  const text = String(seed || "default");
  let hash = 0;
  for (let index = 0; index < text.length; index += 1) {
    hash = (hash << 5) - hash + text.charCodeAt(index);
    hash |= 0;
  }
  return Math.abs(hash);
}

export function defaultCover(seed) {
  return DEFAULT_COVERS[hashSeed(seed) % DEFAULT_COVERS.length];
}

export function cardCover(item, scope = "card", index = 0) {
  if (item?.coverUrl) {
    return item.coverUrl;
  }
  return defaultCover(`${scope}-${item?.id || item?.title || item?.name || index}`);
}
