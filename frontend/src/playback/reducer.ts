import type {
  RiskLevel,
  SandboxEntity,
  SandboxScenario,
  ScenarioEvent,
} from "../scenario/types";

export const PLAYBACK_SPEEDS = [1, 2, 4] as const;

export type PlaybackSpeed = (typeof PLAYBACK_SPEEDS)[number];
export type PlaybackStatus = "paused" | "playing";
export type ScenarioViewState = SandboxScenario["initialState"] & {
  nodeStatuses: Record<string, RiskLevel>;
  edgeStatuses: Record<string, RiskLevel>;
};
export type ScenarioTimelineEvent = SandboxScenario["timeline"][number];

export interface PlaybackState {
  scenario: SandboxScenario;
  currentStep: number;
  status: PlaybackStatus;
  speed: PlaybackSpeed;
  viewState: ScenarioViewState;
}

export type PlaybackAction =
  | { type: "play" }
  | { type: "pause" }
  | { type: "toggle" }
  | { type: "next" }
  | { type: "previous" }
  | { type: "tick" }
  | { type: "seek"; step: number }
  | { type: "set-speed"; speed: PlaybackSpeed }
  | { type: "reset" }
  | { type: "load"; scenario: SandboxScenario };

/**
 * Applies a precomputed scenario delta without mutating either input.
 * Objects merge recursively; arrays and scalar values replace the old value.
 */
export function applyDelta<T>(state: T, delta: unknown): T {
  if (!isRecord(delta) || !isRecord(state)) {
    return clone(delta) as T;
  }

  const result: Record<string, unknown> = { ...state };
  for (const [key, nextValue] of Object.entries(delta)) {
    const currentValue = result[key];
    result[key] = isRecord(currentValue) && isRecord(nextValue)
      ? applyDelta(currentValue, nextValue)
      : clone(nextValue);
  }
  return result as T;
}

/** Rebuilds the exact view state at a timeline boundary. Step 0 is initial state. */
export function replayToStep(
  scenario: SandboxScenario,
  requestedStep: number,
): ScenarioViewState {
  const step = clampStep(requestedStep, scenario.timeline.length);
  let state: ScenarioViewState = {
    ...clone(scenario.initialState),
    nodeStatuses: Object.fromEntries(
      scenario.layout.nodes.map((node) => [node.id, node.status]),
    ),
    edgeStatuses: Object.fromEntries(
      scenario.layout.edges.map((edge) => [edge.id, edge.status]),
    ),
  };

  for (let index = 0; index < step; index += 1) {
    state = applyScenarioEvent(state, scenario.timeline[index]);
  }
  return state;
}

export function createInitialPlaybackState(
  scenario: SandboxScenario,
  speed: PlaybackSpeed = scenario.defaultSpeed,
): PlaybackState {
  return {
    scenario,
    currentStep: 0,
    status: "paused",
    speed,
    viewState: replayToStep(scenario, 0),
  };
}

export function applyScenarioEvent(
  state: ScenarioViewState,
  event: ScenarioEvent,
): ScenarioViewState {
  const entityUpdates = new Map(
    (event.delta.entityUpdates ?? []).map((update) => [update.id, update]),
  );

  return {
    entities: state.entities.map((entity) => {
      const update = entityUpdates.get(entity.id);
      if (!update) return clone(entity);
      return {
        ...entity,
        ...update,
        attributes: update.attributes
          ? { ...entity.attributes, ...update.attributes }
          : { ...entity.attributes },
      } satisfies SandboxEntity;
    }),
    kpis: clone(event.kpis),
    nodeStatuses: applyStatusUpdates(state.nodeStatuses, event.delta.nodeStatuses),
    edgeStatuses: applyStatusUpdates(state.edgeStatuses, event.delta.edgeStatuses),
  };
}

export function playbackReducer(
  state: PlaybackState,
  action: PlaybackAction,
): PlaybackState {
  switch (action.type) {
    case "play":
      return state.currentStep >= state.scenario.timeline.length
        ? { ...state, status: "paused" }
        : { ...state, status: "playing" };
    case "pause":
      return { ...state, status: "paused" };
    case "toggle":
      return playbackReducer(state, { type: state.status === "playing" ? "pause" : "play" });
    case "next":
    case "tick":
      return moveToStep(state, state.currentStep + 1);
    case "previous":
      return moveToStep(state, state.currentStep - 1, true);
    case "seek":
      return moveToStep(state, action.step, true);
    case "set-speed":
      return PLAYBACK_SPEEDS.includes(action.speed)
        ? { ...state, speed: action.speed }
        : state;
    case "reset":
      return {
        ...createInitialPlaybackState(state.scenario, state.speed),
        status: "paused",
      };
    case "load":
      return createInitialPlaybackState(action.scenario, state.speed);
  }
}

export function currentTimelineEvent(state: PlaybackState): ScenarioTimelineEvent | undefined {
  return state.currentStep === 0
    ? undefined
    : state.scenario.timeline[state.currentStep - 1];
}

export function clampStep(step: number, totalSteps: number): number {
  if (!Number.isFinite(step)) return 0;
  return Math.min(totalSteps, Math.max(0, Math.trunc(step)));
}

function moveToStep(
  state: PlaybackState,
  requestedStep: number,
  pause = false,
): PlaybackState {
  const totalSteps = state.scenario.timeline.length;
  const currentStep = clampStep(requestedStep, totalSteps);
  const reachedEnd = currentStep >= totalSteps;

  return {
    ...state,
    currentStep,
    status: pause || reachedEnd ? "paused" : state.status,
    viewState: replayToStep(state.scenario, currentStep),
  };
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function applyStatusUpdates(
  current: Record<string, RiskLevel>,
  updates: Array<{ id: string; status: RiskLevel }> | undefined,
): Record<string, RiskLevel> {
  const next = { ...current };
  for (const update of updates ?? []) next[update.id] = update.status;
  return next;
}

function clone<T>(value: T): T {
  if (typeof structuredClone === "function") return structuredClone(value);
  return JSON.parse(JSON.stringify(value)) as T;
}
