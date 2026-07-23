import type {
  IncorporationFrame,
  IncorporationTrace,
} from "./incorporation";

type RecordValue = Record<string, unknown>;

export type IncorporationStepDefinition = {
  capabilities: string[];
  states: string[];
  entities: string[];
  process: string;
  worldPayloads?: string[];
};

export type IncorporationStepTrace = {
  definition: IncorporationStepDefinition;
  transitions: RecordValue[];
  journal: RecordValue[];
  approvals: RecordValue[];
  decisions: RecordValue[];
  outbox: RecordValue[];
  worldExchanges: RecordValue[];
  processRun?: RecordValue;
  unmatchedCapabilities: string[];
};

export const incorporationStepDefinitions: IncorporationStepDefinition[] = [
  {
    capabilities: ["incorporation.case.open"],
    states: ["incorporation_case_opened"],
    entities: ["incorporation_case"],
    process: "enterprise.incorporation.lifecycle.v1",
  },
  {
    capabilities: [
      "founder.resolution.approve",
      "capital.commitment.record",
      "registration.submit",
    ],
    states: [
      "founder_resolution_approved",
      "capital_commitments_confirmed",
      "registration_submitted",
    ],
    entities: [
      "governance_resolution",
      "capital_commitment",
      "legal_registration",
    ],
    process: "legal.registration.v1",
    worldPayloads: ["registration.review.requested"],
  },
  {
    capabilities: ["registration.observation.commit"],
    states: ["legal_entity_registered"],
    entities: ["legal_registration", "legal_entity"],
    process: "legal.registration.v1",
    worldPayloads: ["registration"],
  },
  {
    capabilities: [
      "bank.account.opening.submit",
      "bank.account.observation.commit",
      "capital.contribution.verify",
    ],
    states: [
      "bank_account_opening_submitted",
      "bank_account_opened",
      "capital_contribution_verified",
    ],
    entities: ["bank_account", "capital_contribution"],
    process: "banking.and.capitalization.v1",
    worldPayloads: ["bank", "capital"],
  },
  {
    capabilities: [
      "organization.establish",
      "executive.appointment.propose",
    ],
    states: ["organization_established"],
    entities: ["organization_position", "executive_appointment"],
    process: "organization.and.appointments.v1",
    worldPayloads: ["appointment"],
  },
  {
    capabilities: [
      "executive.appointment.acceptance.commit",
      "executive.appointment.approve",
      "operating.mandate.grant",
    ],
    states: [
      "executive_appointments_accepted",
      "operating_mandates_activated",
    ],
    entities: ["executive_appointment", "operating_mandate"],
    process: "organization.and.appointments.v1",
    worldPayloads: ["appointment"],
  },
  {
    capabilities: ["initial.budget.approve"],
    states: ["initial_budget_approved"],
    entities: ["budget_envelope"],
    process: "mandate.and.initial.budget.v1",
  },
  {
    capabilities: ["enterprise.readiness.evaluate"],
    states: ["enterprise_operational_ready"],
    entities: [
      "legal_entity",
      "bank_account",
      "operating_mandate",
      "budget_envelope",
    ],
    process: "enterprise.incorporation.lifecycle.v1",
  },
];

const records = (value: unknown): RecordValue[] =>
  Array.isArray(value)
    ? value.filter(
        (item): item is RecordValue =>
          typeof item === "object" && item !== null,
      )
    : [];

const serialized = (value: unknown) => JSON.stringify(value).toLowerCase();

const matchesCapability = (
  item: RecordValue,
  capabilities: string[],
): boolean => {
  const text = serialized(item);
  return capabilities.some((capability) =>
    text.includes(capability.toLowerCase()),
  );
};

const matchesWorld = (
  item: RecordValue,
  definition: IncorporationStepDefinition,
): boolean =>
  matchesCapability(item, definition.capabilities) ||
  (definition.worldPayloads ?? []).some((token) =>
    serialized(item).includes(token.toLowerCase()),
  );

export function buildIncorporationStepTrace(
  frame: IncorporationFrame,
  lifecycle: IncorporationTrace["iaos_lifecycle"],
): IncorporationStepTrace | null {
  const definition = incorporationStepDefinitions[frame.step];
  if (!definition || !lifecycle) return null;
  const processRuns = records(lifecycle.process_runs);
  const processRun =
    processRuns.find(
      (run) =>
        run.process_key === "enterprise.incorporation.lifecycle.v1",
    ) ?? processRuns[0];
  const transitions = records(processRun?.trace).filter((transition) =>
    definition.capabilities.includes(String(transition.capability ?? "")),
  );
  const present = new Set(
    transitions.map((transition) => String(transition.capability ?? "")),
  );
  return {
    definition,
    transitions,
    processRun,
    journal: records(lifecycle.journal).filter((item) =>
      matchesCapability(item, definition.capabilities),
    ),
    approvals: records(lifecycle.approvals).filter((item) =>
      matchesCapability(item, definition.capabilities),
    ),
    decisions: records(lifecycle.decisions).filter((item) =>
      matchesCapability(item, definition.capabilities),
    ),
    outbox: records((lifecycle as RecordValue).outbox).filter((item) =>
      matchesCapability(item, definition.capabilities),
    ),
    worldExchanges: records(lifecycle.world_exchanges).filter((item) =>
      matchesWorld(item, definition),
    ),
    unmatchedCapabilities: definition.capabilities.filter(
      (capability) => !present.has(capability),
    ),
  };
}

export function assertIncorporationStepCoverage(): string[] {
  const all = incorporationStepDefinitions.flatMap(
    (definition) => definition.capabilities,
  );
  const duplicates = all.filter(
    (capability, index) => all.indexOf(capability) !== index,
  );
  return [...new Set(duplicates)];
}
