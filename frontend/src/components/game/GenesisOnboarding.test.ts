import { beforeEach, describe, expect, it } from "vitest";
import {
  clearGenesisOnboardingIdempotencyKey,
  restoreGenesisOnboardingIdempotencyKey,
} from "./GenesisOnboarding";

describe("Genesis onboarding retry identity", () => {
  beforeEach(() => localStorage.clear());

  it("reuses the same idempotency key after a page remount until creation succeeds", () => {
    const firstMount = restoreGenesisOnboardingIdempotencyKey();
    const afterReload = restoreGenesisOnboardingIdempotencyKey();

    expect(afterReload).toBe(firstMount);

    clearGenesisOnboardingIdempotencyKey();
    expect(restoreGenesisOnboardingIdempotencyKey()).not.toBe(firstMount);
  });
});
