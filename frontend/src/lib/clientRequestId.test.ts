import { afterEach, describe, expect, it, vi } from "vitest";
import { createClientRequestId } from "./clientRequestId";

afterEach(() => vi.unstubAllGlobals());

describe("createClientRequestId", () => {
  it("uses randomUUID when the browser exposes it", () => {
    vi.stubGlobal("crypto", { randomUUID: () => "11111111-2222-4333-8444-555555555555" });
    expect(createClientRequestId("site-investigation")).toBe("site-investigation-11111111-2222-4333-8444-555555555555");
  });

  it("uses getRandomValues on LAN HTTP browsers without randomUUID", () => {
    vi.stubGlobal("crypto", {
      getRandomValues: (bytes: Uint8Array) => {
        bytes.fill(0xab);
        return bytes;
      },
    });
    expect(createClientRequestId("site-observation")).toMatch(
      /^site-observation-abababab-abab-4bab-abab-abababababab$/,
    );
  });

  it("still returns unique identifiers in the legacy fallback", () => {
    vi.stubGlobal("crypto", undefined);
    const first = createClientRequestId("site-recommendation");
    const second = createClientRequestId("site-recommendation");
    expect(first).not.toBe(second);
  });
});
