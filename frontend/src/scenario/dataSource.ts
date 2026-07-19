import type { SandboxScenario } from "./types";
import { parseSandboxScenario } from "./validation";

export interface ScenarioDataSource {
  loadScenario(key: string): Promise<SandboxScenario>;
}

/**
 * Deterministic, read-only adapter for bundled preview fixtures.
 * React code depends on ScenarioDataSource rather than importing JSON directly.
 */
export class StaticScenarioDataSource implements ScenarioDataSource {
  constructor(private readonly scenarios: Readonly<Record<string, unknown>>) {}

  async loadScenario(key: string): Promise<SandboxScenario> {
    const source = this.scenarios[key];
    if (source === undefined) {
      throw new Error(`Scenario not found: ${key}`);
    }

    return parseSandboxScenario(source);
  }
}
