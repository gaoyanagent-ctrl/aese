import { useCallback, useEffect, useMemo, useState } from "react";
import dagre from "@dagrejs/dagre";
import {
  Background,
  Controls,
  Handle,
  MarkerType,
  MiniMap,
  Position,
  ReactFlow,
  useNodesState,
  type Edge,
  type Node,
  type NodeProps,
} from "@xyflow/react";
import {
  Activity,
  ArrowLeft,
  ArrowRight,
  BookOpen,
  CheckCircle2,
  Clock3,
  Code2,
  ExternalLink,
  Filter,
  LayoutGrid,
  Network,
  RefreshCw,
  Search,
  X,
} from "lucide-react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import "./SystemAtlas.css";

type Status =
  | "planned"
  | "designed"
  | "building"
  | "validating"
  | "completed"
  | "blocked"
  | "deferred";
export type AtlasEntryRef = {
  label: string;
  app: "iaos" | "aese";
  path: string;
  description?: string;
};
type AtlasNode = {
  id: string;
  key: string;
  project_key: string;
  parent_key?: string;
  title: string;
  subtitle: string;
  kind: string;
  layer: string;
  summary: string;
  target_state: string;
  current_state: string;
  status: Status;
  progress: number;
  health: string;
  owner: string;
  sort_order: number;
  design_refs: string[];
  code_refs: string[];
  evidence_refs: string[];
  entry_refs: AtlasEntryRef[];
  updated_at: string;
};
type AtlasEdge = {
  id: string;
  source: string;
  target: string;
  relation: string;
};
type AtlasUpdate = {
  id: string;
  node_key: string;
  summary: string;
  detail: string;
  occurred_at: string;
  created_by: string;
};
type AtlasData = {
  summary: {
    node_count: number;
    overall_progress: number;
    by_status: Record<string, number>;
    last_updated_at?: string;
  };
  nodes: AtlasNode[];
  edges: AtlasEdge[];
  recent_updates: AtlasUpdate[];
};
type AtlasDocument = { ref: string; title: string; content: string };
const labels: Record<Status, string> = {
  planned: "待规划",
  designed: "已设计",
  building: "建设中",
  validating: "验证中",
  completed: "已完成",
  blocked: "受阻",
  deferred: "后置",
};
const colors: Record<Status, string> = {
  planned: "#64748b",
  designed: "#7c3aed",
  building: "#168aad",
  validating: "#d97706",
  completed: "#059669",
  blocked: "#dc2626",
  deferred: "#78716c",
};
const relationLabels: Record<string, string> = {
  contains: "包含",
  depends_on: "依赖",
  integrates_with: "集成",
  feeds: "供给",
  validates: "验证",
};
const NODE_WIDTH = 245,
  NODE_HEIGHT = 112;

function Card({ data, selected }: NodeProps<Node<AtlasNode>>) {
  return (
    <div
      className={`aese-atlas-node ${selected ? "selected" : ""}`}
      style={{ borderLeftColor: colors[data.status] }}
    >
      <Handle type="target" position={Position.Top} />
      <div>
        <span>{data.layer}</span>
        <i style={{ background: colors[data.status] }} />
      </div>
      <strong>{data.title}</strong>
      <small>{data.subtitle}</small>
      <div className="aese-atlas-meter">
        <i
          style={{
            width: `${data.progress}%`,
            background: colors[data.status],
          }}
        />
      </div>
      <footer>
        <b>{data.progress}%</b>
        <span>{labels[data.status]}</span>
      </footer>
      <Handle type="source" position={Position.Bottom} />
    </div>
  );
}
const nodeTypes = { atlas: Card };

async function atlasToken(forceRefresh = false) {
  const existing = localStorage.getItem("iaos_token");
  if (existing && !forceRefresh) return existing;
  const tenant = localStorage.getItem("aese_iaos_tenant_id") ?? "tenant-hctm";
  const response = await fetch(
    `/api/v1/dev/token?tenant_id=${encodeURIComponent(tenant)}&roles=admin`,
  );
  if (!response.ok) throw new Error("请先通过联动中心建立 IAOS 身份");
  const body = (await response.json()) as { token?: string };
  if (!body.token) throw new Error("IAOS 未返回可用身份");
  localStorage.setItem("iaos_token", body.token);
  return body.token;
}

export async function atlasFetch(input: RequestInfo | URL, init: RequestInit = {}) {
  const request = async (token: string) => {
    const headers = new Headers(init.headers);
    headers.set("Authorization", `Bearer ${token}`);
    return fetch(input, { ...init, headers });
  };
  const response = await request(await atlasToken());
  if (response.status !== 401) return response;
  localStorage.removeItem("iaos_token");
  return request(await atlasToken(true));
}

function autoLayout(nodes: AtlasNode[], edges: AtlasEdge[]): Node<AtlasNode>[] {
  const graph = new dagre.graphlib.Graph().setDefaultEdgeLabel(() => ({}));
  graph.setGraph({
    rankdir: "TB",
    nodesep: 45,
    ranksep: 80,
    marginx: 35,
    marginy: 35,
  });
  nodes.forEach((node) =>
    graph.setNode(node.key, { width: NODE_WIDTH, height: NODE_HEIGHT }),
  );
  edges.forEach((edge) => graph.setEdge(edge.source, edge.target));
  dagre.layout(graph);
  return nodes.map((data) => {
    const point = graph.node(data.key);
    return {
      id: data.key,
      type: "atlas",
      position: { x: point.x - NODE_WIDTH / 2, y: point.y - NODE_HEIGHT / 2 },
      data,
      draggable: true,
    };
  });
}

function Reader({
  document,
  onClose,
}: {
  document: AtlasDocument;
  onClose: () => void;
}) {
  useEffect(() => {
    const close = (event: KeyboardEvent) => event.key === "Escape" && onClose();
    window.addEventListener("keydown", close);
    return () => window.removeEventListener("keydown", close);
  }, [onClose]);
  return (
    <div
      className="aese-atlas-modal"
      onMouseDown={(event) => event.target === event.currentTarget && onClose()}
    >
      <section
        role="dialog"
        aria-modal="true"
        aria-labelledby="aese-atlas-doc-title"
      >
        <header>
          <div>
            <span>{document.ref}</span>
            <h2 id="aese-atlas-doc-title">{document.title}</h2>
          </div>
          <button onClick={onClose} title="关闭文档">
            <X />
          </button>
        </header>
        <article>
          <ReactMarkdown remarkPlugins={[remarkGfm]}>
            {document.content}
          </ReactMarkdown>
        </article>
      </section>
    </div>
  );
}

export function SystemAtlas({
  onExit,
  onNavigate,
}: {
  onExit: () => void;
  onNavigate: (entry: AtlasEntryRef) => void;
}) {
  const [data, setData] = useState<AtlasData | null>(null),
    [selected, setSelected] = useState<AtlasNode | null>(null),
    [error, setError] = useState(""),
    [loading, setLoading] = useState(true),
    [query, setQuery] = useState(""),
    [status, setStatus] = useState<Status | "all">("all");
  const [nodes, setNodes, onNodesChange] = useNodesState<Node<AtlasNode>>([]),
    [document, setDocument] = useState<AtlasDocument | null>(null),
    [documentLoading, setDocumentLoading] = useState("");
  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const response = await atlasFetch("/api/v1/system-atlas?view=aese");
      if (!response.ok) {
        const body = await response.json().catch(() => ({}));
        throw new Error(body.error ?? `HTTP ${response.status}`);
      }
      const next = (await response.json()) as AtlasData;
      next.nodes = next.nodes.map((node) => ({
        ...node,
        design_refs: node.design_refs ?? [],
        code_refs: node.code_refs ?? [],
        evidence_refs: node.evidence_refs ?? [],
        entry_refs: node.entry_refs ?? [],
      }));
      setData(next);
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "全景数据载入失败");
    } finally {
      setLoading(false);
    }
  }, []);
  useEffect(() => {
    void load();
  }, [load]);
  const visibleNodes = useMemo(
    () =>
      (data?.nodes ?? []).filter(
        (node) =>
          (status === "all" || node.status === status) &&
          (!query ||
            `${node.title}${node.subtitle}${node.summary}`
              .toLowerCase()
              .includes(query.toLowerCase())),
      ),
    [data, status, query],
  );
  const visibleKeys = useMemo(
      () => new Set(visibleNodes.map((node) => node.key)),
      [visibleNodes],
    ),
    visibleEdges = useMemo(
      () =>
        (data?.edges ?? []).filter(
          (edge) =>
            visibleKeys.has(edge.source) && visibleKeys.has(edge.target),
        ),
      [data, visibleKeys],
    );
  const applyLayout = useCallback(
    () => setNodes(autoLayout(visibleNodes, visibleEdges)),
    [setNodes, visibleNodes, visibleEdges],
  );
  useEffect(() => {
    applyLayout();
  }, [applyLayout]);
  const relatedKeys = useMemo(() => {
    const keys = new Set<string>();
    if (selected) {
      keys.add(selected.key);
      visibleEdges.forEach((edge) => {
        if (edge.source === selected.key) keys.add(edge.target);
        if (edge.target === selected.key) keys.add(edge.source);
      });
    }
    return keys;
  }, [selected, visibleEdges]);
  const displayNodes = useMemo(
    () =>
      nodes.map((node) => ({
        ...node,
        selected: node.id === selected?.key,
        style: {
          opacity: !selected || relatedKeys.has(node.id) ? 1 : 0.18,
          transition: "opacity 160ms",
        },
      })),
    [nodes, selected, relatedKeys],
  );
  const edges = useMemo<Edge[]>(
    () =>
      visibleEdges.map((edge) => {
        const active =
          !!selected &&
          (edge.source === selected.key || edge.target === selected.key);
        return {
          id: edge.id,
          source: edge.source,
          target: edge.target,
          label: relationLabels[edge.relation] ?? edge.relation,
          animated: active,
          markerEnd: {
            type: MarkerType.ArrowClosed,
            color: active ? "#0f766e" : "#94a3b8",
          },
          style: {
            stroke: active ? "#0f766e" : "#94a3b8",
            strokeWidth: active ? 3 : 1.3,
            opacity: selected && !active ? 0.12 : 1,
          },
          labelStyle: {
            fill: active ? "#0f766e" : "#64748b",
            fontSize: 10,
            fontWeight: active ? 700 : 400,
          },
        };
      }),
    [visibleEdges, selected],
  );
  const relationships = useMemo(
    () =>
      selected
        ? visibleEdges
            .filter(
              (edge) =>
                edge.source === selected.key || edge.target === selected.key,
            )
            .map((edge) => ({
              edge,
              node: data?.nodes.find(
                (node) =>
                  node.key ===
                  (edge.source === selected.key ? edge.target : edge.source),
              ),
              outbound: edge.source === selected.key,
            }))
            .filter((item) => item.node)
        : [],
    [selected, visibleEdges, data],
  );
  const updates = (data?.recent_updates ?? []).filter(
    (update) => !selected || update.node_key === selected.key,
  );
  const readDocument = useCallback(
    async (ref: string) => {
      if (!selected) return;
      setDocumentLoading(ref);
      setError("");
      try {
        const response = await atlasFetch(
          `/api/v1/system-atlas/document?node_key=${encodeURIComponent(selected.key)}&ref=${encodeURIComponent(ref)}`,
        );
        const body = await response.json();
        if (!response.ok)
          throw new Error(body.error ?? `HTTP ${response.status}`);
        setDocument(body as AtlasDocument);
      } catch (reason) {
        setError(reason instanceof Error ? reason.message : "文档读取失败");
      } finally {
        setDocumentLoading("");
      }
    },
    [selected],
  );
  const openEntry = (entry: AtlasEntryRef) => {
    if (entry.app === "aese") {
      onNavigate(entry);
      return;
    }
    const base =
      import.meta.env.VITE_IAOS_URL ??
      `${window.location.protocol}//${window.location.hostname}:3000`;
    window.open(
      `${base.replace(/\/$/, "")}/${entry.path}`,
      "_blank",
      "noopener,noreferrer",
    );
  };
  return (
    <main className="aese-atlas">
      <header className="aese-atlas-header">
        <button
          className="aese-atlas-back"
          onClick={onExit}
          title="返回企业沙盘"
        >
          <ArrowLeft />
        </button>
        <div>
          <span>AESE · ENTERPRISE SIMULATION ANATOMY</span>
          <h1>智能企业仿真完成体</h1>
        </div>
        <div className="aese-atlas-kpis">
          <div>
            <small>总体完成</small>
            <b>{data?.summary.overall_progress ?? "-"}%</b>
          </div>
          <div>
            <small>目标构件</small>
            <b>{data?.summary.node_count ?? "-"}</b>
          </div>
          <div>
            <small>已完成</small>
            <b>{data?.summary.by_status.completed ?? 0}</b>
          </div>
          <div>
            <small>建设 / 验证</small>
            <b>
              {(data?.summary.by_status.building ?? 0) +
                (data?.summary.by_status.validating ?? 0)}
            </b>
          </div>
        </div>
        <button
          className="aese-atlas-refresh"
          onClick={() => void load()}
          title="刷新"
        >
          <RefreshCw className={loading ? "spin" : ""} />
        </button>
      </header>
      <div className="aese-atlas-tools">
        <div>
          <Search />
          <input
            placeholder="搜索模型、场景、Agent 或实验"
            value={query}
            onChange={(event) => setQuery(event.target.value)}
          />
          {query && (
            <button onClick={() => setQuery("")} title="清空">
              <X />
            </button>
          )}
        </div>
        <label>
          <Filter />
          <select
            value={status}
            onChange={(event) =>
              setStatus(event.target.value as Status | "all")
            }
          >
            <option value="all">全部状态</option>
            {Object.entries(labels).map(([key, value]) => (
              <option key={key} value={key}>
                {value}
              </option>
            ))}
          </select>
        </label>
      <button
        className="aese-atlas-layout"
        onClick={applyLayout}
        aria-label="恢复自动布局"
        title="恢复自动布局"
      >
        <LayoutGrid />
      </button>
        <span>
          拖动节点可调整 ·{" "}
          {data?.summary.last_updated_at
            ? new Date(data.summary.last_updated_at).toLocaleString("zh-CN")
            : "等待连接"}
        </span>
      </div>
      {error ? (
        <div className="aese-atlas-notice">
          <Activity />
          <span>{error}</span>
          <button onClick={() => setError("")}>
            <X />
          </button>
        </div>
      ) : null}
      <div className="aese-atlas-body">
        <section className="aese-atlas-canvas">
          <div className="aese-atlas-legend">
            {Object.entries(labels).map(([key, value]) => (
              <span key={key}>
                <i style={{ background: colors[key as Status] }} />
                {value}
              </span>
            ))}
          </div>
          <ReactFlow
            nodes={displayNodes}
            edges={edges}
            nodeTypes={nodeTypes}
            onNodesChange={onNodesChange}
            onNodeClick={(_, node) => setSelected(node.data as AtlasNode)}
            fitView
            fitViewOptions={{ minZoom: 0.4, padding: 0.12 }}
            minZoom={0.25}
            maxZoom={1.7}
            proOptions={{ hideAttribution: true }}
          >
            <Background color="#cbd5e1" gap={22} />
            <Controls position="bottom-left" />
            <MiniMap
              pannable
              zoomable
              nodeColor={(node) => colors[(node.data as AtlasNode).status]}
            />
          </ReactFlow>
        </section>
        <aside className={`aese-atlas-detail ${selected ? "open" : ""}`}>
          {selected ? (
            <>
              <div className="aese-atlas-detail-head">
                <div>
                  <span>
                    {selected.layer} · {selected.owner}
                  </span>
                  <h2>{selected.title}</h2>
                  <p>{selected.subtitle}</p>
                </div>
                <button onClick={() => setSelected(null)} title="关闭">
                  <X />
                </button>
              </div>
              <div className="aese-atlas-progress">
                <b>{selected.progress}%</b>
                <div>
                  <i
                    style={{
                      width: `${selected.progress}%`,
                      background: colors[selected.status],
                    }}
                  />
                </div>
                <em style={{ color: colors[selected.status] }}>
                  {labels[selected.status]}
                </em>
              </div>
              <section>
                <h3>仿真职责</h3>
                <p>{selected.summary}</p>
              </section>
              <section>
                <h3>完成体目标</h3>
                <p>{selected.target_state}</p>
              </section>
              <section>
                <h3>当前状态</h3>
                <p>{selected.current_state}</p>
              </section>
              <section>
                <h3>
                  <Network />
                  相关构件
                </h3>
                {relationships.length ? (
                  <div className="aese-atlas-relations">
                    {relationships.map(({ edge, node, outbound }) => (
                      <button key={edge.id} onClick={() => setSelected(node!)}>
                        <span>
                          {outbound ? "本构件" : "上游"}{" "}
                          {relationLabels[edge.relation] ?? edge.relation}{" "}
                          {outbound ? "下游" : "本构件"}
                        </span>
                        <b>{node!.title}</b>
                        <ArrowRight />
                      </button>
                    ))}
                  </div>
                ) : (
                  <p>当前筛选范围内无直接关系</p>
                )}
              </section>
              <section>
                <h3>
                  <BookOpen />
                  设计文档
                </h3>
                {selected.design_refs.length ? (
                  <div className="aese-atlas-actions">
                    {selected.design_refs.map((ref) => (
                      <button
                        key={ref}
                        disabled={documentLoading === ref}
                        onClick={() => void readDocument(ref)}
                      >
                        <BookOpen />
                        <code>{ref}</code>
                        <ExternalLink />
                      </button>
                    ))}
                  </div>
                ) : (
                  <p>尚未登记设计文档</p>
                )}
              </section>
              <section>
                <h3>
                  <ArrowRight />
                  功能入口
                </h3>
                {selected.entry_refs.length ? (
                  <div className="aese-atlas-entries">
                    {selected.entry_refs.map((entry) => (
                      <button
                        key={`${entry.app}:${entry.path}`}
                        onClick={() => openEntry(entry)}
                      >
                        <div>
                          <b>{entry.label}</b>
                          <span>
                            {entry.description ??
                              `${entry.app.toUpperCase()} · ${entry.path}`}
                          </span>
                        </div>
                        <ExternalLink />
                      </button>
                    ))}
                  </div>
                ) : (
                  <p>该构件尚无可操作菜单</p>
                )}
              </section>
              <section>
                <h3>
                  <Code2 />
                  代码位置
                </h3>
                {selected.code_refs.length ? (
                  selected.code_refs.map((ref) => (
                    <div className="aese-atlas-ref" key={ref}>
                      <Code2 />
                      <code>{ref}</code>
                    </div>
                  ))
                ) : (
                  <p>尚未登记代码位置</p>
                )}
              </section>
              <section>
                <h3>验证证据</h3>
                {selected.evidence_refs.length ? (
                  <div className="aese-atlas-actions">
                    {selected.evidence_refs.map((ref) => (
                      <button key={ref} onClick={() => void readDocument(ref)}>
                        <CheckCircle2 />
                        <code>{ref}</code>
                        <ExternalLink />
                      </button>
                    ))}
                  </div>
                ) : (
                  <p>尚未登记验证证据</p>
                )}
              </section>
            </>
          ) : (
            <div className="aese-atlas-empty">
              <Activity />
              <b>点击构件查看详情</b>
              <p>查看直接关系、设计文档、可运行菜单、代码位置和进展。</p>
            </div>
          )}
          <section className="aese-atlas-timeline">
            <h3>
              <Clock3 />
              进展记录
            </h3>
            {updates.slice(0, 8).map((update) => (
              <article key={update.id}>
                <i />
                <div>
                  <b>{update.summary}</b>
                  <p>{update.detail}</p>
                  <span>
                    {new Date(update.occurred_at).toLocaleString("zh-CN")} ·{" "}
                    {update.created_by}
                  </span>
                </div>
              </article>
            ))}
          </section>
        </aside>
      </div>
      <footer className="aese-atlas-footer">
        <CheckCircle2 />
        数据来自 IAOS System Atlas
        数据库；选择节点可查看一跳关系，节点支持拖动调整。
      </footer>
      {document && (
        <Reader document={document} onClose={() => setDocument(null)} />
      )}
    </main>
  );
}
