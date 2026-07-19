import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import type { SandboxEdge, SandboxNode } from '../scenario/types';
import { FactoryCanvas } from './FactoryCanvas';

class ResizeObserverStub implements ResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
}

vi.stubGlobal('ResizeObserver', ResizeObserverStub);

afterEach(cleanup);

const nodes: SandboxNode[] = [
  {
    id: 'supplier-primary',
    businessCode: 'SUP-AL-001',
    label: '主供铝材供应商',
    kind: 'supplier',
    position: { x: 0, y: 120 },
    status: 'watch',
  },
  {
    id: 'raw-warehouse',
    businessCode: 'WH-RM-SZ-01',
    label: '原材仓',
    kind: 'warehouse',
    position: { x: 260, y: 120 },
    status: 'normal',
  },
];

const edges: SandboxEdge[] = [
  {
    id: 'supply-path',
    source: 'supplier-primary',
    target: 'raw-warehouse',
    label: '铝材供应',
    status: 'watch',
  },
];

describe('FactoryCanvas', () => {
  it('renders an explicit empty state when the scenario has no nodes', () => {
    render(<FactoryCanvas nodes={[]} edges={[]} />);

    expect(screen.getByText('暂无工厂布局')).toBeInTheDocument();
    expect(screen.getByText('场景数据中没有可显示的节点。')).toBeInTheDocument();
  });

  it('exposes status text and activates a node with pointer and keyboard input', () => {
    const onNodeSelect = vi.fn();
    const { rerender } = render(
      <FactoryCanvas
        nodes={nodes}
        edges={edges}
        highlightedNodeIds={['supplier-primary']}
        highlightedEdgeIds={['supply-path']}
        onNodeSelect={onNodeSelect}
      />,
    );

    const supplier = screen.getByLabelText('主供铝材供应商，关注，查看详情');
    fireEvent.click(supplier);
    fireEvent.keyDown(supplier, { key: 'Enter' });

    expect(onNodeSelect).toHaveBeenCalledTimes(2);
    expect(onNodeSelect).toHaveBeenLastCalledWith(nodes[0]);
    expect(supplier).toHaveAttribute('aria-pressed', 'false');

    rerender(
      <FactoryCanvas
        nodes={nodes}
        edges={edges}
        selectedNodeId="supplier-primary"
        onNodeSelect={onNodeSelect}
      />,
    );
    expect(
      screen.getByLabelText('主供铝材供应商，关注，查看详情'),
    ).toHaveAttribute('aria-pressed', 'true');
  });
});
