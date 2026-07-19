import { useEffect, useMemo, useReducer } from "react";

import type { SandboxScenario } from "../scenario/types";
import {
  createInitialPlaybackState,
  currentTimelineEvent,
  playbackReducer,
  type PlaybackSpeed,
} from "./reducer";

const BASE_EVENT_INTERVAL_MS = 1_000;

export interface PlaybackControls {
  play: () => void;
  pause: () => void;
  toggle: () => void;
  next: () => void;
  previous: () => void;
  seek: (step: number) => void;
  reset: () => void;
  setSpeed: (speed: PlaybackSpeed) => void;
}

export function usePlayback(scenario: SandboxScenario) {
  const [state, dispatch] = useReducer(
    playbackReducer,
    scenario,
    createInitialPlaybackState,
  );

  useEffect(() => {
    dispatch({ type: "load", scenario });
  }, [scenario]);

  useEffect(() => {
    if (state.status !== "playing") return undefined;

    const intervalId = window.setInterval(() => {
      dispatch({ type: "tick" });
    }, BASE_EVENT_INTERVAL_MS / state.speed);

    return () => window.clearInterval(intervalId);
  }, [state.speed, state.status]);

  const controls = useMemo<PlaybackControls>(() => ({
    play: () => dispatch({ type: "play" }),
    pause: () => dispatch({ type: "pause" }),
    toggle: () => dispatch({ type: "toggle" }),
    next: () => dispatch({ type: "next" }),
    previous: () => dispatch({ type: "previous" }),
    seek: (step) => dispatch({ type: "seek", step }),
    reset: () => dispatch({ type: "reset" }),
    setSpeed: (speed) => dispatch({ type: "set-speed", speed }),
  }), []);

  const currentEvent = useMemo(() => currentTimelineEvent(state), [state]);
  const totalSteps = scenario.timeline.length;

  return {
    ...state,
    currentEvent,
    totalSteps,
    isAtStart: state.currentStep === 0,
    isAtEnd: state.currentStep === totalSteps,
    controls,
  };
}

export type PlaybackController = ReturnType<typeof usePlayback>;
