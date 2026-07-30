import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { approveAndExecuteWorkItem, createIncorporationCase, listGenesisWorkspaces, signInGenesisPlayer } from "./api";

describe("Genesis workspace player session recovery", () => {
  beforeEach(() => {
    localStorage.clear();
    localStorage.setItem("aese_genesis_player_id", "player-local-test");
    localStorage.setItem("aese_genesis_player_token", "expired-player-token");
    localStorage.setItem("iaos_token", "current-iaos-token");
  });

  afterEach(() => vi.restoreAllMocks());

  it("retries the current IAOS session when the persisted player token expired", async () => {
    const fetchMock = vi.spyOn(globalThis, "fetch")
      .mockResolvedValueOnce(new Response(JSON.stringify({
        error: "IAOS player session expired",
        code: "player_session_expired",
      }), { status: 401 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({
        items: [{ workspace_id: "gxw-test", status: "active" }],
      }), { status: 200 }));

    const items = await listGenesisWorkspaces();

    expect(items).toHaveLength(1);
    expect((fetchMock.mock.calls[0][1]?.headers as Record<string, string>).Authorization)
      .toBe("Bearer expired-player-token");
    expect((fetchMock.mock.calls[1][1]?.headers as Record<string, string>).Authorization)
      .toBe("Bearer current-iaos-token");
    expect(localStorage.getItem("aese_genesis_player_token")).toBe("current-iaos-token");
  });

  it("clears an expired player token and reports a login action when no session works", async () => {
    vi.spyOn(globalThis, "fetch")
      .mockResolvedValueOnce(new Response(JSON.stringify({
        error: "IAOS player session expired",
        code: "player_session_expired",
      }), { status: 401 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({
        error: "IAOS player session expired",
        code: "player_session_expired",
      }), { status: 401 }));

    await expect(listGenesisWorkspaces()).rejects.toThrow("IAOS 登录已过期");
    expect(localStorage.getItem("aese_genesis_player_token")).toBeNull();
  });
});

describe("createIncorporationCase", () => {
  beforeEach(() => {
    localStorage.clear();
    localStorage.setItem("iaos_token", "stale-admin-token");
    localStorage.setItem("aese_iaos_tenant_id", "tenant-gx-test");
    localStorage.setItem("aese_genesis_workspace_id", "gxw-test");
    localStorage.setItem("aese_genesis_player_id", "player-local-test");
  });

  afterEach(() => vi.restoreAllMocks());

  it("refreshes the owner founder session once after the legacy principal mismatch", async () => {
    const fetchMock = vi.spyOn(globalThis, "fetch")
      .mockResolvedValueOnce(new Response(
        JSON.stringify({ error: "governance rejected: authenticated principal does not match acting subject or tenant access" }),
        { status: 422 },
      ))
      .mockResolvedValueOnce(new Response(JSON.stringify({
        workspace_id: "gxw-test",
        tenant_id: "tenant-gx-test",
        tenant_token: "founder-token",
      }), { status: 200 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ duplicate: false }), { status: 201 }));

    await createIncorporationCase({
      case_code: "INC-GX-TEST",
      case_name: "测试企业设立案",
      proposed_company_name: "测试企业有限公司",
      registered_address: "江苏省苏州市测试大道1号",
      business_scope: "工业热管理系统研发与制造",
    });

    expect(fetchMock).toHaveBeenCalledTimes(3);
    expect(fetchMock.mock.calls[1][0]).toBe("/api/aese/v1/genesis/workspaces/gxw-test/session");
    expect((fetchMock.mock.calls[2][1]?.headers as Record<string, string>).Authorization).toBe("Bearer founder-token");
    expect(localStorage.getItem("iaos_token")).toBe("founder-token");
  });
});

describe("local Genesis player login", () => {
  it("claims the existing browser player for the first username", () => {
    localStorage.clear();
    localStorage.setItem("aese_genesis_player_id", "player-local-existing");

    const session = signInGenesisPlayer("founder-principal");

    expect(session.player_id).toBe("player-local-existing");
    expect(localStorage.getItem("aese_genesis_username")).toBe("founder-principal");
  });
});

describe("approveAndExecuteWorkItem", () => {
  beforeEach(() => {
    localStorage.clear();
    localStorage.setItem("iaos_token", "founder-token");
    localStorage.setItem("aese_iaos_tenant_id", "tenant-gx-test");
  });

  afterEach(() => vi.restoreAllMocks());

  it("keeps approval-only business notes out of the strict IAOS execute command", async () => {
    const fetchMock = vi.spyOn(globalThis, "fetch")
      .mockResolvedValueOnce(new Response(JSON.stringify({ approval: { id: "approval-g3" } }), { status: 201 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ status: "approved" }), { status: 200 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ status: "committed" }), { status: 201 }));

    await approveAndExecuteWorkItem(
      "INC-GX-TEST", 8, "bank.account.opening.submit", "G3",
      { business_note: "申请开户银行：江海商业银行；资料：营业执照、公章" },
    );

    const gateBody = JSON.parse(String(fetchMock.mock.calls[0][1]?.body));
    const executeBody = JSON.parse(String(fetchMock.mock.calls[2][1]?.body));
    expect(gateBody.intent.business_note).toContain("江海商业银行");
    expect(executeBody).toEqual({ correlation_id: "corr-game-INC-GX-TEST-8" });
  });

  it("does not decide an already-approved request again when retrying execute", async () => {
    const fetchMock = vi.spyOn(globalThis, "fetch")
      .mockResolvedValueOnce(new Response(JSON.stringify({
        approval: { id: "approval-g3", status: "approved" },
      }), { status: 200 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ status: "committed" }), { status: 201 }));

    await approveAndExecuteWorkItem(
      "INC-GX-TEST", 8, "bank.account.opening.submit", "G3",
      { business_note: "申请开户银行：江海商业银行" },
    );

    expect(fetchMock).toHaveBeenCalledTimes(2);
    expect(String(fetchMock.mock.calls[1][0])).toContain("/work-items/8/execute");
  });

  it("uses a fresh governed correlation when a monetary approval is corrected", async () => {
    const fetchMock = vi.spyOn(globalThis, "fetch")
      .mockResolvedValueOnce(new Response(JSON.stringify({ approval: { id: "approval-g4-new" } }), { status: 201 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ status: "approved" }), { status: 200 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ status: "committed" }), { status: 201 }));

    await approveAndExecuteWorkItem(
      "INC-GX-TEST", 10, "capital.contribution.verify", "G4",
      { amount_minor: 100_000_000, currency: "CNY" },
    );

    const gateBody = JSON.parse(String(fetchMock.mock.calls[0][1]?.body));
    const executeBody = JSON.parse(String(fetchMock.mock.calls[2][1]?.body));
    expect(gateBody.correlation_id).toBe("corr-game-INC-GX-TEST-10-100000000");
    expect(executeBody).toEqual({
      amount_minor: 100_000_000,
      currency: "CNY",
      correlation_id: "corr-game-INC-GX-TEST-10-100000000",
    });
  });
});
