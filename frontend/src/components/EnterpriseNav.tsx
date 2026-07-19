import { Boxes, Building2, Factory, MapPin, PackageCheck, Warehouse } from 'lucide-react';

interface EnterpriseNavProps {
  selectedId: string | null;
  onSelect: (id: string) => void;
}

const items = [
  { id: 'HCTM-GROUP', label: '华辰热管理集团', meta: '集团', icon: Building2 },
  { id: 'PLT-SZ', label: '苏州制造基地', meta: '基地 · 运行中', icon: MapPin },
  { id: 'LINE-BCP-A', label: '电池冷却板 A 线', meta: '产线 · 重点', icon: Factory },
  { id: 'WH-SZ-RM', label: '原材料仓', meta: '仓储', icon: Warehouse },
  { id: 'WH-SZ-FG', label: '成品仓', meta: '仓储', icon: PackageCheck },
];

export function EnterpriseNav({ selectedId, onSelect }: EnterpriseNavProps) {
  return (
    <nav className="enterprise-nav" aria-label="虚拟企业对象">
      <div className="panel-heading">
        <span className="eyebrow">ENTERPRISE</span>
        <h2>企业结构</h2>
      </div>
      <div className="nav-tree">
        {items.map(({ id, label, meta, icon: Icon }, index) => (
          <button
            key={id}
            className={`nav-item level-${Math.min(index, 2)} ${selectedId === id ? 'selected' : ''}`}
            onClick={() => onSelect(id)}
            aria-pressed={selectedId === id}
          >
            <Icon aria-hidden="true" />
            <span><strong>{label}</strong><small>{meta}</small></span>
          </button>
        ))}
      </div>
      <div className="line-summary" aria-label="A 线运行摘要">
        <Boxes aria-hidden="true" />
        <div><strong>12 个关键区域</strong><span>8 道核心工艺 · 4 个物流节点</span></div>
      </div>
    </nav>
  );
}
