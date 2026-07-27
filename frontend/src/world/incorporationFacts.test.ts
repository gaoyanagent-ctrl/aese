import { describe, expect, it } from "vitest";
import {
  resolveIncorporationFacts,
  type IncorporationFrame,
} from "./incorporation";

const frame = {
  company: {
    code: "BANK-HCTM-SZ-OPERATING-01",
    owner: "HCTM-SZ-MFG",
    status: "not_opened",
    balance: { value: "0.00", currency: "CNY", scale: 2 },
  },
  capital_committed: { value: "30000000.00", currency: "CNY", scale: 2 },
  capital_paid: { value: "0.00", currency: "CNY", scale: 2 },
  budget: {
    code: "BUDGET-HCTM-SZ-Y1",
    status: "not_submitted",
    amount: { value: "15000000.00", currency: "CNY", scale: 2 },
    owner: "HCTM-SZ-MFG",
  },
} as IncorporationFrame;

const lifecycle = (state: Record<string, unknown>) => ({
  case_code: "INC-001",
  state,
  journal: [],
  approvals: [],
  outbox: [],
  world_exchanges: [],
  process_runs: [],
  decisions: [],
  runtime_artifact: {},
  lineage: {},
});

describe("incorporation display fact lineage", () => {
  it("labels the pack target as a scenario assumption without an IAOS case", () => {
    expect(resolveIncorporationFacts(frame).legalEntity).toEqual({
      value: "HCTM-SZ-MFG（场景目标）",
      source: "world_baseline",
      sourceLabel: "AESE 华辰场景假设",
    });
  });

  it("shows the proposed name before IAOS assigns a legal entity code", () => {
    expect(
      resolveIncorporationFacts(
        frame,
        lifecycle({ proposed_company_name: "苏州测试制造有限公司" }),
      ).legalEntity,
    ).toMatchObject({
      value: "苏州测试制造有限公司（待登记编码）",
      source: "pending",
    });
  });

  it("prefers IAOS committed identifiers and monetary facts", () => {
    const facts = resolveIncorporationFacts(
      frame,
      lifecycle({
        legal_entity_code: "LEGAL-INC-001",
        bank_account_code: "BANK-INC-001",
        commitment_minor: 1234500,
        contribution_minor: 1000000,
        budget_authorized_minor: 800000,
      }),
    );
    expect(facts.legalEntity.value).toBe("LEGAL-INC-001");
    expect(facts.companyAccount.value).toBe("BANK-INC-001");
    expect(facts.capitalCommitted.value).toBe("12345.00");
    expect(facts.capitalPaid.value).toBe("10000.00");
    expect(facts.budget.value).toBe("8000.00");
    expect(facts.budget.source).toBe("iaos_committed");
  });
});
