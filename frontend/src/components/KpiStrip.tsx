import { Activity, Boxes, CircleGauge, PackageOpen, Truck } from 'lucide-react';
import type { KpiMetric, KpiSnapshot } from '../scenario/types';

interface KpiStripProps { kpis: KpiSnapshot }

const definitions: Array<{ key: keyof KpiSnapshot; label: string; icon: typeof Activity }> = [
  { key: 'orderDemand', label: '订单需求', icon: PackageOpen },
  { key: 'availableFinishedGoods', label: '可用成品', icon: Boxes },
  { key: 'materialShortageRisk', label: '缺料风险', icon: Activity },
  { key: 'capacityRisk', label: '产能风险', icon: CircleGauge },
  { key: 'deliveryRisk', label: '交付风险', icon: Truck },
];

function formatMetric(metric: KpiMetric) {
  if (metric.unit === '%') return `${metric.value}%`;
  return `${new Intl.NumberFormat('zh-CN').format(metric.value)}${metric.unit ? ` ${metric.unit}` : ''}`;
}

export function KpiStrip({ kpis }: KpiStripProps) {
  return (
    <section className="kpi-strip" aria-label="关键经营指标">
      {definitions.map(({ key, label, icon: Icon }) => {
        const metric = kpis[key];
        return (
          <article className={`kpi-card status-${metric.risk}`} key={key}>
            <div className="kpi-label"><Icon aria-hidden="true" /><span>{label}</span></div>
            <strong>{formatMetric(metric)}</strong>
            <span className="kpi-state">{metric.risk === 'normal' ? '正常' : metric.risk === 'watch' ? '需关注' : '高风险'} · {metric.trend === 'up' ? '上升' : metric.trend === 'down' ? '下降' : '持平'}</span>
          </article>
        );
      })}
    </section>
  );
}
