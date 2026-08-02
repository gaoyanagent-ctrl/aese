let fallbackSequence = 0;

function bytesToUUID(bytes: Uint8Array): string {
  bytes[6] = (bytes[6] & 0x0f) | 0x40;
  bytes[8] = (bytes[8] & 0x3f) | 0x80;
  const hex = Array.from(bytes, (value) => value.toString(16).padStart(2, "0")).join("");
  return `${hex.slice(0, 8)}-${hex.slice(8, 12)}-${hex.slice(12, 16)}-${hex.slice(16, 20)}-${hex.slice(20)}`;
}

/**
 * Generates a collision-resistant browser request identifier on both HTTPS
 * and LAN HTTP origins. randomUUID is secure-context-only in some browsers,
 * while getRandomValues remains available on those HTTP origins.
 */
export function createClientRequestId(prefix: string): string {
  const browserCrypto = globalThis.crypto as Partial<Crypto> | undefined;
  if (typeof browserCrypto?.randomUUID === "function") {
    return `${prefix}-${browserCrypto.randomUUID()}`;
  }
  if (typeof browserCrypto?.getRandomValues === "function") {
    const bytes = new Uint8Array(16);
    browserCrypto.getRandomValues(bytes);
    return `${prefix}-${bytesToUUID(bytes)}`;
  }
  fallbackSequence += 1;
  const timestamp = Date.now().toString(36);
  const random = Math.random().toString(16).slice(2).padEnd(13, "0").slice(0, 13);
  return `${prefix}-${timestamp}-${fallbackSequence.toString(36)}-${random}`;
}
