import {
  ChevronLeft,
  ChevronRight,
  Pause,
  Play,
  PlugZap,
  Radio,
  RefreshCw,
  RotateCcw,
} from 'lucide-react';

interface ControlBarProps {
  scenarioName: string;
  currentTime: string;
  step: number;
  totalSteps: number;
  playing: boolean;
  speed: 1 | 2 | 4;
  onTogglePlay: () => void;
  onPrevious: () => void;
  onNext: () => void;
  onReset: () => void;
  onSpeedChange: (speed: 1 | 2 | 4) => void;
  mode?: 'preview' | 'live';
  onModeChange?: (mode: 'preview' | 'live') => void;
  sourceStatus?: string;
  onRefresh?: () => void;
  onReconnect?: () => void;
  onOpenIntegration?: () => void;
}

export function ControlBar({
  scenarioName,
  currentTime,
  step,
  totalSteps,
  playing,
  speed,
  onTogglePlay,
  onPrevious,
  onNext,
  onReset,
  onSpeedChange,
  mode = 'preview',
  onModeChange,
  sourceStatus,
  onRefresh,
  onReconnect,
  onOpenIntegration,
}: ControlBarProps) {
  return (
    <header className="control-bar">
      <div className="brand-block">
        <span className="brand-mark" aria-hidden="true">AE</span>
        <div>
          <div className="eyebrow">AESE · ENTERPRISE SIMULATION</div>
          <h1>{scenarioName}</h1>
        </div>
      </div>

      <div className="source-clock" aria-label="仿真状态">
        <div className="mode-switch" aria-label="数据模式">
          <button className={mode === 'preview' ? 'active' : ''} onClick={() => onModeChange?.('preview')}>Preview</button>
          <button className={mode === 'live' ? 'active' : ''} onClick={() => onModeChange?.('live')}>Live</button>
        </div>
        <span className={`source-pill ${mode === 'live' ? 'source-live' : ''}`}><span className="source-dot" />{mode === 'live' ? (sourceStatus ?? 'LIVE') : 'PREVIEW'}</span>
        <span className="clock-value">{currentTime}</span>
        <span className="step-value">事件 {step}/{totalSteps}</span>
      </div>

      <div className="playback-controls" aria-label="故事播放控制">
        <button className="integration-button" onClick={onOpenIntegration} aria-label="打开 AESE 与 IAOS 联动中心"><PlugZap aria-hidden="true" /><span>联动中心</span></button>
        {mode === 'live' ? <>
          <button className="play-button" onClick={onRefresh} aria-label="刷新在线快照"><RefreshCw aria-hidden="true" /><span>刷新</span></button>
          <button className="icon-button" onClick={onReconnect} aria-label="重新连接在线事件"><Radio aria-hidden="true" /></button>
        </> : <>
        <button className="icon-button" onClick={onPrevious} disabled={step === 0} aria-label="上一个事件">
          <ChevronLeft aria-hidden="true" />
        </button>
        <button className="play-button" onClick={onTogglePlay} disabled={step === totalSteps} aria-label={playing ? '暂停故事' : '播放故事'}>
          {playing ? <Pause aria-hidden="true" /> : <Play aria-hidden="true" />}
          <span>{playing ? '暂停' : '播放'}</span>
        </button>
        <button className="icon-button" onClick={onNext} disabled={step === totalSteps} aria-label="下一个事件">
          <ChevronRight aria-hidden="true" />
        </button>
        <label className="speed-control">
          <span className="sr-only">播放速度</span>
          <select value={speed} onChange={(event) => onSpeedChange(Number(event.target.value) as 1 | 2 | 4)}>
            <option value={1}>1×</option>
            <option value={2}>2×</option>
            <option value={4}>4×</option>
          </select>
        </label>
        <button className="icon-button reset-button" onClick={onReset} disabled={step === 0} aria-label="重置故事">
          <RotateCcw aria-hidden="true" />
        </button>
        </>}
      </div>
    </header>
  );
}
