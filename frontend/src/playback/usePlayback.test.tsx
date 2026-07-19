import { act, renderHook } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { createPlaybackFixture } from "./testFixture";
import { usePlayback } from "./usePlayback";

describe("usePlayback", () => {
  beforeEach(() => vi.useFakeTimers());
  afterEach(() => {
    vi.clearAllTimers();
    vi.useRealTimers();
  });

  it("plays, pauses, steps, seeks and resets", () => {
    const scenario = createPlaybackFixture();
    const { result } = renderHook(() => usePlayback(scenario));

    act(() => result.current.controls.next());
    expect(result.current.currentStep).toBe(1);
    expect(result.current.currentEvent?.id).toBe("event-1");

    act(() => result.current.controls.previous());
    expect(result.current.currentStep).toBe(0);

    act(() => result.current.controls.seek(2));
    expect(result.current.currentStep).toBe(2);

    act(() => result.current.controls.reset());
    expect(result.current.currentStep).toBe(0);
    expect(result.current.status).toBe("paused");
  });

  it("honours 1x/2x/4x and stops its interval at the end", () => {
    const scenario = createPlaybackFixture(2);
    const { result } = renderHook(() => usePlayback(scenario));

    act(() => {
      result.current.controls.setSpeed(4);
      result.current.controls.play();
    });
    act(() => vi.advanceTimersByTime(250));
    expect(result.current.currentStep).toBe(1);

    act(() => vi.advanceTimersByTime(250));
    expect(result.current.currentStep).toBe(2);
    expect(result.current.status).toBe("paused");
    expect(vi.getTimerCount()).toBe(0);
  });

  it("cleans the timer on pause and unmount", () => {
    const scenario = createPlaybackFixture();
    const { result, unmount } = renderHook(() => usePlayback(scenario));
    act(() => result.current.controls.play());
    expect(vi.getTimerCount()).toBe(1);

    act(() => result.current.controls.pause());
    expect(vi.getTimerCount()).toBe(0);

    act(() => result.current.controls.play());
    expect(vi.getTimerCount()).toBe(1);
    unmount();
    expect(vi.getTimerCount()).toBe(0);
  });
});
