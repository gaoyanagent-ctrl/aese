import { Box, X } from 'lucide-react';
import type { SandboxEntity } from '../scenario/types';

interface EntityDrawerProps { entity: SandboxEntity | null; onClose: () => void }

export function EntityDrawer({ entity, onClose }: EntityDrawerProps) {
  if (!entity) return null;
  return (
    <div className="drawer-backdrop" onMouseDown={(event) => event.target === event.currentTarget && onClose()}>
      <aside className="entity-drawer" role="dialog" aria-modal="true" aria-labelledby="entity-title">
        <button className="icon-button drawer-close" onClick={onClose} aria-label="关闭对象详情"><X aria-hidden="true" /></button>
        <div className="drawer-icon"><Box aria-hidden="true" /></div>
        <span className="eyebrow">{entity.type}</span>
        <h2 id="entity-title">{entity.name}</h2>
        <div className="entity-code">{entity.businessCode}</div>
        <div className={`entity-risk status-${entity.risk}`}>状态：{entity.status} · {entity.risk === 'normal' ? '正常' : entity.risk === 'watch' ? '关注' : '严重'}</div>
        <dl className="attribute-list">
          {Object.entries(entity.attributes).map(([key, value]) => <div key={key}><dt>{key}</dt><dd>{String(value ?? '—')}</dd></div>)}
        </dl>
      </aside>
    </div>
  );
}
