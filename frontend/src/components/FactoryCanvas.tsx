import {
  Background,
  BackgroundVariant,
  Controls,
  Handle,
  MarkerType,
  MiniMap,
  Position,
  ReactFlow,
  type Edge,
  type Node,
  type NodeProps,
} from '@xyflow/react';
import {
  Boxes,
  Building2,
  CheckCircle2,
  CircleAlert,
  Cog,
  Factory,
  PackageCheck,
  ScanLine,
  ShieldCheck,
  Truck,
  TriangleAlert,
  Warehouse,
  Waves,
  Wrench,
  type LucideIcon,
} from 'lucide-react';
import { useMemo } from 'react';
import type { KeyboardEvent } from 'react';
import type {
  RiskLevel,
  SandboxEdge,
  SandboxNode,
} from '../scenario/types';
import '@xyflow/react/dist/style.css';
import './FactoryCanvas.css';

export interface FactoryCanvasProps {
  nodes: readonly SandboxNode[];
  edges: readonly SandboxEdge[];
  selectedNodeId?: string | null;
  highlightedNodeIds?: readonly string[];
  highlightedEdgeIds?: readonly string[];
  onNodeSelect?: (node: SandboxNode) => void;
  className?: string;
  ariaLabel?: string;
}

type FactoryNodeData = {
  source: SandboxNode;
  label: string;
  code: string;
  kind: string;
  status: RiskLevel;
  statusLabel: string;
  highlighted: boolean;
  selected: boolean;
  onActivate?: (node: SandboxNode) => void;
};

type FactoryFlowNode = Node<FactoryNodeData, 'factory'>;

const nodeTypes = { factory: FactoryNode };

const iconByKind: Record<string, LucideIcon> = {
  supplier: Factory,
  raw_material_warehouse: Warehouse,
  warehouse: Warehouse,
  process: Cog,
  forming: Wrench,
  machining: Cog,
  welding: Wrench,
  cleaning: Waves,
  leak_test: ScanLine,
  assembly: Boxes,
  packaging: PackageCheck,
  finished_goods_warehouse: Warehouse,
  shipping: Truck,
  customer: Building2,
  quality: ShieldCheck,
};

const iconByCodePrefix: Record<string, LucideIcon> = {
  FRM: Wrench,
  CNC: Cog,
  WLD: Wrench,
  CLN: Waves,
  LKT: ScanLine,
  ASM: Boxes,
  PKG: PackageCheck,
  IQC: ShieldCheck,
};

function iconForNode(kind: string, code: string): LucideIcon {
  const prefix = code.split('-')[0];
  return iconByCodePrefix[prefix] ?? iconByKind[kind] ?? Boxes;
}

const statusLabel: Record<RiskLevel, string> = {
  normal: '正常',
  watch: '关注',
  critical: '严重',
};

const statusIcon: Record<RiskLevel, LucideIcon> = {
  normal: CheckCircle2,
  watch: TriangleAlert,
  critical: CircleAlert,
};

function FactoryNode({ data }: NodeProps<FactoryFlowNode>) {
  const KindIcon = iconForNode(data.kind, data.code);
  const StatusIcon = statusIcon[data.status];

  return (
    <div
      className={[
        'factory-node',
        `factory-node--${data.status}`,
        data.highlighted ? 'factory-node--highlighted' : '',
        data.selected ? 'factory-node--selected' : '',
      ]
        .filter(Boolean)
        .join(' ')}
      data-status={data.status}
      role="button"
      tabIndex={0}
      aria-pressed={data.selected}
      aria-label={`${data.label}，${data.statusLabel}，查看详情`}
      onClick={() => data.onActivate?.(data.source)}
      onKeyDown={(event: KeyboardEvent<HTMLDivElement>) => {
        if (event.key === 'Enter' || event.key === ' ') {
          event.preventDefault();
          data.onActivate?.(data.source);
        }
      }}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="factory-node__handle"
        isConnectable={false}
      />
      <span className="factory-node__icon" aria-hidden="true">
        <KindIcon size={20} strokeWidth={1.8} />
      </span>
      <span className="factory-node__content">
        <strong title={data.label}>{data.label}</strong>
        <span className="factory-node__code" title={data.code}>
          {data.code}
        </span>
      </span>
      <span className="factory-node__status">
        <StatusIcon size={14} aria-hidden="true" />
        <span>{data.statusLabel}</span>
      </span>
      <Handle
        type="source"
        position={Position.Right}
        className="factory-node__handle"
        isConnectable={false}
      />
    </div>
  );
}

function edgeStatus(edge: SandboxEdge): RiskLevel {
  return edge.status ?? 'normal';
}

export function FactoryCanvas({
  nodes,
  edges,
  selectedNodeId = null,
  highlightedNodeIds = [],
  highlightedEdgeIds = [],
  onNodeSelect,
  className,
  ariaLabel = '苏州制造基地电池冷却板 A 线工艺与物流画布',
}: FactoryCanvasProps) {
  const highlightedNodes = useMemo(
    () => new Set(highlightedNodeIds),
    [highlightedNodeIds],
  );
  const highlightedEdges = useMemo(
    () => new Set(highlightedEdgeIds),
    [highlightedEdgeIds],
  );
  const sourceNodes = useMemo(
    () => new Map(nodes.map((node) => [node.id, node])),
    [nodes],
  );

  const flowNodes = useMemo<FactoryFlowNode[]>(
    () =>
      nodes.map((node) => {
        const status = node.status ?? 'normal';
        const selected = node.id === selectedNodeId;

        return {
          id: node.id,
          type: 'factory',
          position: node.position,
          initialWidth: 200,
          initialHeight: 76,
          draggable: false,
          selectable: true,
          focusable: false,
          selected,
          ariaLabel: `${node.label}，${statusLabel[status]}，按回车查看详情`,
          data: {
            source: node,
            label: node.label,
            code: node.businessCode,
            kind: node.kind,
            status,
            statusLabel: statusLabel[status],
            highlighted: highlightedNodes.has(node.id),
            selected,
            onActivate: onNodeSelect,
          },
        };
      }),
    [highlightedNodes, nodes, onNodeSelect, selectedNodeId],
  );

  const flowEdges = useMemo<Edge[]>(
    () =>
      edges.map((edge) => {
        const status = edgeStatus(edge);
        const highlighted = highlightedEdges.has(edge.id);

        return {
          id: edge.id,
          source: edge.source,
          target: edge.target,
          type: 'smoothstep',
          animated: highlighted,
          focusable: true,
          ariaLabel:
            edge.label ??
            `${sourceNodes.get(edge.source)?.label ?? edge.source} 到 ${sourceNodes.get(edge.target)?.label ?? edge.target}`,
          className: [
            'factory-edge',
            `factory-edge--${status}`,
            highlighted ? 'factory-edge--highlighted' : '',
          ]
            .filter(Boolean)
            .join(' '),
          markerEnd: {
            type: MarkerType.ArrowClosed,
            width: 14,
            height: 14,
            color: highlighted
              ? 'var(--factory-canvas-selected)'
              : status === 'critical'
                ? 'var(--factory-canvas-critical)'
                : status === 'watch'
                  ? 'var(--factory-canvas-attention)'
                  : '#718197',
          },
          label: edge.label,
          labelShowBg: true,
          labelBgPadding: [5, 3],
          labelBgBorderRadius: 4,
        };
      }),
    [edges, highlightedEdges, sourceNodes],
  );

  const rootClassName = ['factory-canvas', className].filter(Boolean).join(' ');

  if (nodes.length === 0) {
    return (
      <section className={`${rootClassName} factory-canvas--empty`} aria-label={ariaLabel}>
        <Factory aria-hidden="true" />
        <strong>暂无工厂布局</strong>
        <span>场景数据中没有可显示的节点。</span>
      </section>
    );
  }

  return (
    <section className={rootClassName} aria-label={ariaLabel}>
      <ReactFlow<FactoryFlowNode, Edge>
        nodes={flowNodes}
        edges={flowEdges}
        nodeTypes={nodeTypes}
        fitView
        fitViewOptions={{ padding: 0.14, minZoom: 0.2, maxZoom: 1.1 }}
        minZoom={0.2}
        maxZoom={1.8}
        nodesDraggable={false}
        nodesConnectable={false}
        elementsSelectable
        panOnDrag
        zoomOnPinch
        zoomOnScroll={false}
        preventScrolling={false}
        proOptions={{ hideAttribution: true }}
        colorMode="dark"
      >
        <Background
          variant={BackgroundVariant.Dots}
          gap={20}
          size={1.2}
          color="var(--factory-canvas-grid, #334155)"
        />
        <MiniMap
          className="factory-canvas__minimap"
          pannable
          zoomable
          nodeStrokeWidth={2}
          nodeColor={(node) => {
            if (node.data.selected) return 'var(--factory-canvas-selected)';
            if (node.data.status === 'critical') return 'var(--factory-canvas-critical)';
            if (node.data.status === 'watch') return 'var(--factory-canvas-attention)';
            return 'var(--factory-canvas-normal)';
          }}
          ariaLabel="A 线布局导航图"
        />
        <Controls
          className="factory-canvas__controls"
          position="bottom-left"
          showInteractive={false}
          aria-label="画布缩放与适配控件"
        />
      </ReactFlow>
      <div className="factory-canvas__legend" aria-label="运行状态图例">
        <span data-status="normal"><CheckCircle2 aria-hidden="true" />正常</span>
        <span data-status="attention"><TriangleAlert aria-hidden="true" />关注</span>
        <span data-status="critical"><CircleAlert aria-hidden="true" />严重</span>
        <span data-status="selected"><ScanLine aria-hidden="true" />已选中</span>
      </div>
    </section>
  );
}
