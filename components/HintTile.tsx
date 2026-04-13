'use client'

import { HintResult } from '@/lib/game'

interface HintTileProps {
  label: string
  value: string | number
  result: HintResult
  direction?: 'up' | 'down' | 'exact'
  delay?: number
  isEmpty?: boolean
}

const config = {
  green: {
    bg: '#e8f8ed',
    border: '#34c759',
    text: '#1c6b32',
    indicator: '#34c759',
    label: '#34c759',
  },
  'half-green': {
    bg: 'linear-gradient(90deg, #e8f8ed 50%, #f7f7f7 50%)',
    border: '#34c759',
    text: '#1c6b32',
    indicator: '#34c759',
    label: '#34c759',
  },
  yellow: {
    bg: '#fff5e6',
    border: '#ff9f0a',
    text: '#7a4500',
    indicator: '#ff9f0a',
    label: '#ff9f0a',
  },
  red: {
    bg: '#fff0ef',
    border: '#ff3b30',
    text: '#7a1a15',
    indicator: '#ff3b30',
    label: '#ff3b30',
  },
}

const arrows = { up: '↑', down: '↓', exact: '' }

export default function HintTile({ label, value, result, direction, delay = 0, isEmpty }: HintTileProps) {
  if (isEmpty) {
    return (
      <div
        style={{
          background: '#ffffff',
          border: '1.5px solid #e0e0e5',
          borderRadius: '10px',
          aspectRatio: '1 / 1.15',
          width: '100%',
        }}
      />
    )
  }

  const c = config[result]
  const isHalf = result === 'half-green'

  return (
    <div
      className="tile-reveal relative flex flex-col items-center justify-center overflow-hidden"
      style={{
        animationDelay: `${delay}ms`,
        background: isHalf ? undefined : c.bg,
        backgroundImage: isHalf ? c.bg : undefined,
        border: `1.5px solid ${c.border}`,
        borderRadius: '10px',
        aspectRatio: '1 / 1.15',
        width: '100%',
        padding: '6px 4px',
      }}
    >
      {/* Half-green divider */}
      {isHalf && (
        <div
          className="absolute inset-y-0 left-1/2 w-px"
          style={{ background: '#34c75940' }}
        />
      )}

      {/* Category label */}
      <div
        className="absolute top-1.5 left-0 right-0 text-center"
        style={{
          fontSize: '8px',
          fontWeight: 600,
          color: c.label,
          letterSpacing: '0.04em',
          fontFamily: '"DM Mono", monospace',
          textTransform: 'uppercase',
          opacity: 0.75,
        }}
      >
        {label}
      </div>

      {/* Value */}
      <div
        className="text-center leading-tight"
        style={{
          fontSize: 'clamp(9px, 1.5vw, 13px)',
          fontWeight: 600,
          color: c.text,
          fontFamily: '"DM Sans", sans-serif',
          wordBreak: 'break-word',
          maxWidth: '100%',
          paddingTop: '10px',
        }}
      >
        {value}
      </div>

      {/* Direction arrow */}
      {direction && direction !== 'exact' && (
        <div
          className="absolute bottom-1.5"
          style={{
            fontSize: '10px',
            color: c.indicator,
            fontWeight: 700,
            lineHeight: 1,
          }}
        >
          {arrows[direction]}
        </div>
      )}
    </div>
  )
}
